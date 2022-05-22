package api

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrOrderPaymentNotFound      = errors.New("order payment not found")
	ErrOrderPaymentNotAuthorized = errors.New("payment not authorized")
	ErrOrderPaymentNotCompleted  = errors.New("payment not completed")
)

type PaymentAPI interface {
	AuthorizeOrder(orderID uuid.UUID, totalAmount int) error
	CompleteTransaction(orderID uuid.UUID) error
	CancelOrder(orderID uuid.UUID) error
	RefundOrder(orderID uuid.UUID) error
}
