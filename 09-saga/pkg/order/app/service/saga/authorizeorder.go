package saga

import (
	"errors"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/saga"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/domain"
)

type authorizeOrderOperation struct {
	paymentAPI api.PaymentAPI
	order      *domain.Order
	logger     log.Logger
}

func (op *authorizeOrderOperation) Name() string {
	return "authorizeOrder"
}

func (op *authorizeOrderOperation) Do() error {
	err := op.paymentAPI.AuthorizeOrder(op.order.ID, op.order.TotalAmount)
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.order.ID,
		}).WithError(err).Error("failed to authorize order")
		return err
	}
	return nil
}

func (op *authorizeOrderOperation) Undo() error {
	err := op.paymentAPI.CancelOrder(op.order.ID)
	if errors.Is(err, api.ErrOrderPaymentNotAuthorized) {
		return nil
	}
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.order.ID,
		}).WithError(err).Error("failed to cancel order")
		return err
	}
	return nil
}

func NewAuthorizeOrderOperation(
	paymentAPI api.PaymentAPI,
	order *domain.Order,
	logger log.Logger,
) saga.Operation {
	return &authorizeOrderOperation{
		paymentAPI: paymentAPI,
		order:      order,
		logger:     logger,
	}
}
