package deliveryapi

import (
	"fmt"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
)

type apiClient struct {
	serviceURL string
}

func (a *apiClient) ScheduleDelivery() error {
	//TODO implement me
	fmt.Println("schedule delivery")
	return nil
}

func New(serviceURL string) api.DeliveryAPI {
	return &apiClient{serviceURL: serviceURL}
}
