package domain

import (
	"github.com/google/uuid"
)

type OrderStatus int

const (
	OrderStatusCreated OrderStatus = iota
)

type OrderItem struct {
	ID        uuid.UUID
	ItemPrice int
	Quantity  int
}

type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Items       []OrderItem
	Status      OrderStatus
	TotalAmount int
}

type OrderRepository interface {
	NextID() uuid.UUID
	Store(order *Order) error
}
