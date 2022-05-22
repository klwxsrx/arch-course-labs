package api

type DeliveryAPI interface {
	ScheduleDelivery() error
	DeleteDeliverySchedule() error
}
