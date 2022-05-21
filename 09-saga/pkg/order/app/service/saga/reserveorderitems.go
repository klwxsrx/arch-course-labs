package saga

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/saga"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/app/service/api"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/order/domain"
)

type reserveOrderItemsOperation struct {
	stockAPI api.StockAPI
	order    *domain.Order
	logger   log.Logger
}

func (op *reserveOrderItemsOperation) Name() string {
	return "reserveOrderItems"
}

func (op *reserveOrderItemsOperation) Do() error {
	itemIDs := make([]uuid.UUID, 0, len(op.order.Items))
	for _, item := range op.order.Items {
		itemIDs = append(itemIDs, item.ID)
	}

	err := op.stockAPI.ReserveOrderItems(op.order.ID, itemIDs)
	if errors.Is(err, api.ErrOrderItemsOutOfStock) {
		return fmt.Errorf("failed to reserve items: %w", err)
	}
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.order.ID,
		}).WithError(err).Error("failed to reserve order items")
		return err
	}
	return nil
}

func (op *reserveOrderItemsOperation) Undo() error {
	err := op.stockAPI.RemoveOrderItemsReservation(op.order.ID)
	if err != nil {
		op.logger.With(log.Fields{
			"orderID": op.order.ID,
		}).WithError(err).Error("failed to remove order item reservation")
		return err
	}
	return nil
}

func NewReserveOrderItemsOperation(
	stockAPI api.StockAPI,
	order *domain.Order,
	logger log.Logger,
) saga.Operation {
	return &reserveOrderItemsOperation{
		stockAPI: stockAPI,
		order:    order,
		logger:   logger,
	}
}
