package api

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrOrderPaymentRejected = errors.New("order payment rejected")
)

type PaymentAPI interface {
	AuthorizeOrder(orderID uuid.UUID, totalAmount int) error
	CompleteTransaction(orderID uuid.UUID) error
	CancelOrder(orderID uuid.UUID) error
	RefundOrder(orderID uuid.UUID) error
}
