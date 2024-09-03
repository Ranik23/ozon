package usecase

import (
	"ozon1/internal/entity/item"
	"ozon1/internal/repository"
	"time"
)


type UseCase interface {
	AcceptFromDelivery(order_id, receiver_id int, deadline time.Time) error 
	ReturnToDelivery(order_id int) error
	ReturnToClient(order_id... int) (*[]item.Item, error)
	ReturnAllTheOrders(receiver_id int, others... int) (*[]item.Item, error) 
	AcceptReturning(receiver_id, order_id int) error
	GetTheReturnList() ([]item.Item, error)
}

type UserOperator struct {
	repo *repository.Repository
}

func NewUserOperator(repository *repository.Repository) *UserOperator {
	return &UserOperator{
		repo : repository,
	}
}

func (r *UserOperator) AcceptFromDelivery(order_id, receiver_id int, deadline time.Time) error {
	return r.repo.AcceptFromDelivery(order_id, receiver_id, deadline)
}

func (r *UserOperator) ReturnToDelivery(order_id int) error {
	return r.repo.ReturnToDelivery(order_id)
}

func (r *UserOperator) ReturnToClient(order_id... int) (*[]item.Item, error) {
	return r.repo.ReturnToClient(order_id...)
}

func (r* UserOperator) ReturnAllTheOrders(receiver_id int , limit int) (*[]item.Item, error) {
	return r.repo.GetOrders(receiver_id, limit)
}

func (r* UserOperator) AcceptReturning(receiver_id, order_id int) error {
	return r.repo.AcceptReturning(receiver_id, order_id)
}

func (r *UserOperator) GetTheReturnList() (*[]item.Item, error) {
	return r.repo.GetTheReturnList()
}







