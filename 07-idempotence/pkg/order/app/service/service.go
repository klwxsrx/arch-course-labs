package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/app/idempotence"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/app/persistence"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/domain"
)

var (
	ErrOrderAlreadyCreated = errors.New("order with key is already created")
	ErrEmptyOrder          = errors.New("empty or completely free order")
)

type OrderService struct {
	ufw    persistence.UnitOfWork
	logger log.Logger
}

func (s *OrderService) Create(
	idempotenceKey string,
	userID uuid.UUID,
	items []domain.OrderItem,
) error {
	err := s.ufw.Execute(func(p persistence.PersistentProvider) error {
		err := p.IdempotenceKeyStore().StoreUnique(idempotenceKey)
		if errors.Is(err, idempotence.ErrKeyAlreadyExists) {
			return ErrOrderAlreadyCreated
		}
		if err != nil {
			return err
		}

		totalAmount := s.calculateTotalAmount(items)
		if totalAmount == 0 {
			return ErrEmptyOrder
		}

		order := &domain.Order{
			ID:          p.OrderRepository().NextID(),
			UserID:      userID,
			Items:       items,
			Status:      domain.OrderStatusCreated,
			TotalAmount: totalAmount,
		}

		return p.OrderRepository().Store(order)
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create order")
	}
	return err
}

func (s *OrderService) calculateTotalAmount(items []domain.OrderItem) int {
	var result int
	for _, item := range items {
		result += item.ItemPrice * item.Quantity
	}
	return result
}

func NewOrderService(ufw persistence.UnitOfWork, logger log.Logger) *OrderService {
	return &OrderService{ufw: ufw, logger: logger}
}
