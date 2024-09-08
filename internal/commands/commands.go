package commands

import (
	"fmt"
	"ozon1/internal/usecase"
	"ozon1/internal/utils"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)


func AcceptOrderFromDelivery(uc *usecase.UserOperator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "accept-order",
		Short: "Accept order from courier",
		Run: func(cmd *cobra.Command, args []string) {

			orderID, _ := cmd.Flags().GetString("order-id")

			userID, _ := cmd.Flags().GetString("user-id")

			expirationStr, _ := cmd.Flags().GetString("expiration")

			expiration, err := time.Parse("2006-01-02", expirationStr)

			if err != nil {
				fmt.Println("Invalid date format")
				return
			}

			order_id, _ := strconv.Atoi(orderID)
			user_id, _ := strconv.Atoi(userID)

			if err := uc.AcceptFromDelivery(order_id, user_id, expiration); err != nil {
				fmt.Println("Error accepting order:", err)
			} else {
				fmt.Println("Order accepted")
			}
		},
	}

	cmd.Flags().String("order-id", "", "ID of the order")
	cmd.Flags().String("user-id", "", "ID of the user")
	cmd.Flags().String("expiration", "", "Expiration date (YYYY-MM-DD)")
	cmd.MarkFlagRequired("order-id")
	cmd.MarkFlagRequired("user-id")
	cmd.MarkFlagRequired("expiration")

	return cmd
}

func ReturnOrderToDelivery(uc *usecase.UserOperator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "return-order-to-delivery",
		Short: "Return order to delivery",
		Run: func(cmd *cobra.Command, args []string) {
			orderID, _ := cmd.Flags().GetString("order-id")

			order_id, _ := strconv.Atoi(orderID)

			if err := uc.ReturnToDelivery(order_id); err != nil {
				fmt.Println("Error returning order:", err)
			} else {
				fmt.Println("Order returned to courier")
			}
		},
	}

	cmd.Flags().String("order-id", "", "ID of the order")
	cmd.MarkFlagRequired("order-id")

	return cmd
}

func ReturnOrderToClient(uc *usecase.UserOperator) *cobra.Command{
	cmd := &cobra.Command{
		Use:   "return-order-client",
		Short: "Return order to client",
		Run: func(cmd *cobra.Command, args []string) {
			orderID, _ := cmd.Flags().GetString("order-id")

			order_id, _ := strconv.Atoi(orderID)

			items, err := uc.ReturnToClient(order_id)

			if err != nil {
				fmt.Println("Error returning order:", err)
			} else {
				fmt.Println("Order returned to client")
				fmt.Println(items)
			}
		},
	}

	cmd.Flags().String("order-id", "", "ID of the order")
	cmd.MarkFlagRequired("order-id")

	return cmd
}

func ReturnAllOrders(uc *usecase.UserOperator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "return-all-orders",
		Short: "Return all orders",
		Run : func(cmd *cobra.Command, args []string) {
			receiverID, _ := cmd.Flags().GetString("receiver-id")
			limit, _ := cmd.Flags().GetString("limit")

			receiver_id, _ := strconv.Atoi(receiverID)
			Limit, _ := strconv.Atoi(limit)

			items, err := uc.ReturnAllTheOrders(receiver_id, Limit)

			if err != nil {
				fmt.Println("Error returning all orders", err)
			} else {
				utils.Paginate(*items)
			}
		},
	}

	cmd.Flags().String("receiver-id", "", "ID of the receiver")
	cmd.Flags().String("limit", "", "limit")
	cmd.MarkFlagRequired("receiver-id")
	cmd.MarkFlagRequired("limit")

	return cmd
}

func AcceptReturn(uc *usecase.UserOperator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "accept-return",
		Short: "Accept the return",
		Run: func(cmd *cobra.Command, args []string) {
			receiverID, _ := cmd.Flags().GetString("receiver-id")
			orderID, _ := cmd.Flags().GetString("order-id")
			receiver_id, _ := strconv.Atoi(receiverID)
			order_id, _ := strconv.Atoi(orderID)

			if err := uc.AcceptReturning(receiver_id, order_id); err != nil{
				fmt.Println("Error accepting the return", err)
			} else {
				fmt.Println("accepted the return")
			}
		},
	}

	cmd.Flags().String("receiver-id", "", "ID of the receiver")
	cmd.Flags().String("order-id", "", "ID of the order")
	cmd.MarkFlagRequired("receiver-id")
	cmd.MarkFlagRequired("order-id")

	return cmd
}

func GetReturnList(uc *usecase.UserOperator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-return-list",
		Short: "Get the return list",
		Run: func(cmd *cobra.Command, args []string) {

			items, err := uc.GetTheReturnList();
			if err != nil {
				fmt.Println("error return all the items")
			} else {
				fmt.Println(*items)
			}
		},
	}
	return cmd
}




