package saga

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/saga"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
)

type completeTransactionOperation struct {
	paymentAPI api.PaymentAPI
	orderID    uuid.UUID
	logger     log.Logger
}

func (op *completeTransactionOperation) Name() string {
	return "completeTransaction"
}

func (op *completeTransactionOperation) Do() error {
	err := op.paymentAPI.CompleteTransaction(op.orderID)
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.orderID,
		}).WithError(err).Error("failed to complete transaction")
		return err
	}
	return nil
}

func (op *completeTransactionOperation) Undo() error {
	err := op.paymentAPI.RefundOrder(op.orderID)
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.orderID,
		}).WithError(err).Error("failed to refund order")
		return err
	}
	return nil
}

func NewCompleteTransactionOperation(
	paymentAPI api.PaymentAPI,
	orderID uuid.UUID,
	logger log.Logger,
) saga.Operation {
	return &completeTransactionOperation{
		paymentAPI: paymentAPI,
		orderID:    orderID,
		logger:     logger,
	}
}
