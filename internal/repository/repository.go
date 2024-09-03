package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"ozon1/internal/entity/item"
	"time"
)

type UserRepository interface {
	AcceptFromDelivery(order_id, receiver_id int, deadline time.Time) error
	ReturnToDelivery(order_id int) error
	ReturnToClient(order_id... int) (*[]item.Item, error)
	GetOrders(receiver_id int, limit int) (*[]item.Item, error)
	AcceptReturning(receiver_id, order_id int) error
	GetTheReturnList() (*[]item.Item, error)
}

type Repository struct {
	log                 *slog.Logger
	orders_path         string
	returned_orders_path string
}

func NewRepository(log *slog.Logger, PathToOrders, PathToReturnedOrders string) *Repository {
	return &Repository{
		log:                 log,
		orders_path:         PathToOrders,
		returned_orders_path: PathToReturnedOrders,
	}
}

func (r *Repository) readOrders(filePath string) ([]item.Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		r.log.Error("failed to open the file", slog.String("err", err.Error()), slog.String("file", filePath))
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		r.log.Error("failed to read the file", slog.String("err", err.Error()), slog.String("file", filePath))
		return nil, err
	}

	var orders []item.Item
	if err = json.Unmarshal(bytes, &orders); err != nil && len(bytes) != 0 {
		r.log.Error("failed to unmarshal", slog.String("err", err.Error()))
		return nil, err
	}

	return orders, nil
}

func (r *Repository) writeOrders(filePath string, orders []item.Item) error {
	file, err := os.Create(filePath)
	if err != nil {
		r.log.Error("failed to create the file", slog.String("err", err.Error()))
		return err
	}
	defer file.Close()

	data, err := json.Marshal(orders)
	if err != nil {
		r.log.Error("failed to marshal orders", slog.String("err", err.Error()))
		return err
	}

	if _, err = file.Write(data); err != nil {
		r.log.Error("failed to write to file", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func remove[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func (r *Repository) AcceptFromDelivery(order_id, receiver_id int, deadline time.Time) error {
	if deadline.Before(time.Now()) {
		return fmt.Errorf("the date has expired")
	}

	orders, err := r.readOrders(r.orders_path)

	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.Id == order_id && order.Receiver_Id == receiver_id {
			return fmt.Errorf("twice added")
		}
	}

	newOrder := item.Item{
		Id:              order_id,
		Receiver_Id:     receiver_id,
		Expiration_Date: deadline,
		Received:        false,
		Status:          true,
	}

	orders = append(orders, newOrder)

	return r.writeOrders(r.orders_path, orders)
}

func (r *Repository) ReturnToDelivery(order_id int) error {

	orders, err := r.readOrders(r.orders_path)

	if err != nil {
		return err
	}

	var found bool = false

	var ordersReturned []item.Item

	for index, order := range orders {
		if order.Id == order_id  && !order.Received {
			orders = remove(orders, index)
			ordersReturned = append(ordersReturned, order)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("not found")
	}
	
	if err := r.writeOrders(r.returned_orders_path, ordersReturned); err != nil {
		return fmt.Errorf("error writing to returned orders")
	}
	return r.writeOrders(r.orders_path, orders)
}

func (r *Repository) ReturnToClient(order_ids ...int) (*[]item.Item, error) {

	orders, err := r.readOrders(r.orders_path)

	if err != nil {
		return nil, err
	}

	var clientID int

	for _, orderID := range order_ids {
		found := false
		for i, order := range orders {
			if order.Id == orderID {
				if clientID == 0 {
					clientID = order.Receiver_Id
				} else if order.Receiver_Id != clientID {
					return nil, fmt.Errorf("not all orders belong to the same receiver")
				}

				if order.Expiration_Date.Before(time.Now()) {
					return nil, fmt.Errorf("the order with ID %d has expired and cannot be returned", orderID)
				}

				orders = remove(orders, i)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("order with ID %d not found", orderID)
		}
	}

	return &orders, r.writeOrders(r.orders_path, orders)
}

func (r *Repository) GetOrders(receiver_id int, limit int) (*[]item.Item, error) {

	orders, err := r.readOrders(r.orders_path)

	if err != nil {
		return nil, err
	}

	var filteredOrders []item.Item

	for _, order := range orders {
		if order.Receiver_Id == receiver_id && order.Status && order.Expiration_Date.After(time.Now()) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	if limit > 0 && len(filteredOrders) > limit {
		filteredOrders = filteredOrders[:limit]
	}

	return &filteredOrders, nil
}


func (r *Repository) AcceptReturning(receiver_id, order_id int) error {

	orders, err := r.readOrders(r.orders_path)

	if err != nil {
		return err
	}

	var updatedOrders []item.Item
	var found bool

	for _, order := range orders {
		if order.Id == order_id {
			found = true
			if order.Receiver_Id != receiver_id {
				return fmt.Errorf("order ID %d does not belong to receiver ID %d", order_id, receiver_id)
			}

			if time.Since(order.Expiration_Date) > 48*time.Hour {
				return fmt.Errorf("return period for order ID %d has expired", order_id)
			}

			if !order.Status {
				return fmt.Errorf("order ID %d was not issued from our PВЗ", order_id)
			}
		}
		updatedOrders = append(updatedOrders, order)
	}

	if !found {
		return fmt.Errorf("order ID %d not found", order_id)
	}

	return r.writeOrders(r.returned_orders_path, updatedOrders)
}

func (r *Repository) GetTheReturnList() (*[]item.Item, error) {

	orders, err := r.readOrders(r.orders_path)
	
	if err != nil {
		return nil, err
	}

	return &orders, nil
}
