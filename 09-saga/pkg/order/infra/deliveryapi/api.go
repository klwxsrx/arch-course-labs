package deliveryapi

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
)

type apiClient struct {
	serviceURL string
}

func (a *apiClient) ScheduleDelivery(orderID uuid.UUID, addressID uuid.UUID) error {
	//TODO implement me
	fmt.Println("schedule delivery")
	return nil
}

func (a *apiClient) DeleteDeliverySchedule(orderID uuid.UUID) error {
	//TODO implement me
	fmt.Println("delete delivery schedule")
	return nil
}

func New(serviceURL string) api.DeliveryAPI {
	return &apiClient{serviceURL: serviceURL}
}
