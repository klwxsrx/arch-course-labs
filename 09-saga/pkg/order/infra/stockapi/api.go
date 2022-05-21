package stockapi

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
)

type apiClient struct {
	serviceURL string
}

func (a *apiClient) ReserveOrderItems(orderID uuid.UUID, itemIDs []uuid.UUID) error {
	//TODO implement me
	fmt.Println("reserve order items")
	return nil
}

func (a *apiClient) RemoveOrderItemsReservation(orderID uuid.UUID) error {
	//TODO implement me
	fmt.Println("remove order item reservation")
	return nil
}

func New(serviceURL string) api.StockAPI {
	return &apiClient{serviceURL: serviceURL}
}
