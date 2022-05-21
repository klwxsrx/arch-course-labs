package api

import (
	"errors"
	"github.com/google/uuid"
)

var ErrOrderItemsOutOfStock = errors.New("items out of stock")

type StockAPI interface {
	ReserveOrderItems(orderID uuid.UUID, itemIDs []uuid.UUID) error
	RemoveOrderItemsReservation(orderID uuid.UUID) error
}
