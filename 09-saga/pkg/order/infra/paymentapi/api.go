package paymentapi

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
)

type apiClient struct {
	serviceURL string
}

func (a *apiClient) AuthorizeOrder(orderID uuid.UUID, totalAmount int) error {
	//TODO implement me
	fmt.Println("authorize order")
	return nil
}

func (a *apiClient) CompleteTransaction(orderID uuid.UUID) error {
	//TODO implement me
	fmt.Println("complete transaction")
	return nil
}

func (a *apiClient) CancelOrder(orderID uuid.UUID) error {
	//TODO implement me
	fmt.Println("cancel order")
	return nil
}

func New(serviceURL string) api.PaymentAPI {
	return &apiClient{serviceURL: serviceURL}
}
