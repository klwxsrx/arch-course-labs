package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/payment/app/persistence"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/payment/domain"
)

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentNotAuthorized = errors.New("payment not authorized")
)

type PaymentService struct {
	ufw    persistence.UnitOfWork
	logger log.Logger
}

func (s *PaymentService) CreatePayment(orderID uuid.UUID, totalAmount int) error {
	err := s.ufw.Execute(func(p persistence.PersistentProvider) error {
		// TODO: authorize payment from payment gateway
		payment, err := p.PaymentRepository().GetByID(orderID)
		if err != nil && !errors.Is(err, domain.ErrPaymentNotFound) {
			return err
		}
		if err == nil {
			return nil
		}

		payment = &domain.Payment{
			OrderID:     orderID,
			TotalAmount: totalAmount,
			Status:      domain.PaymentStatusAuthorized,
		}

		return p.PaymentRepository().Store(payment)
	})
	if err != nil {
		s.logger.WithError(err).With(log.Fields{
			"orderID":     orderID,
			"totalAmount": totalAmount,
		}).Error("failed to create payment")
		return err
	}
	return nil
}

func (s *PaymentService) CompletePayment(orderID uuid.UUID) error {
	err := s.ufw.Execute(func(p persistence.PersistentProvider) error {
		// TODO: complete payment from payment gateway

		payment, err := p.PaymentRepository().GetByID(orderID)
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return ErrPaymentNotFound
		}
		if err != nil {
			return err
		}

		if payment.Status == domain.PaymentStatusCompleted {
			return nil
		}
		if payment.Status != domain.PaymentStatusAuthorized {
			return ErrPaymentNotAuthorized
		}

		payment.Status = domain.PaymentStatusCompleted
		return p.PaymentRepository().Store(payment)
	})
	if errors.Is(err, ErrPaymentNotFound) || errors.Is(err, ErrPaymentNotAuthorized) {
		return nil
	}
	if err != nil {
		s.logger.WithError(err).With(log.Fields{
			"orderID": orderID,
		}).Error("failed to complete payment")
		return err
	}
	return nil
}

func (s *PaymentService) CancelPayment(orderID uuid.UUID) error {
	err := s.ufw.Execute(func(p persistence.PersistentProvider) error {
		// TODO: cancel payment from payment gateway

		payment, err := p.PaymentRepository().GetByID(orderID)
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return ErrPaymentNotFound
		}
		if err != nil {
			return err
		}

		if payment.Status == domain.PaymentStatusCancelled {
			return nil
		}
		if payment.Status != domain.PaymentStatusAuthorized {
			return ErrPaymentNotAuthorized
		}

		payment.Status = domain.PaymentStatusCancelled
		return p.PaymentRepository().Store(payment)
	})
	if errors.Is(err, ErrPaymentNotFound) || errors.Is(err, ErrPaymentNotAuthorized) {
		return nil
	}
	if err != nil {
		s.logger.WithError(err).With(log.Fields{
			"orderID": orderID,
		}).Error("failed to cancel payment")
		return err
	}
	return nil
}

func NewPaymentService(ufw persistence.UnitOfWork, logger log.Logger) *PaymentService {
	return &PaymentService{ufw: ufw, logger: logger}
}
