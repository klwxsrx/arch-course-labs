package persistence

import (
	"github.com/klwxsrx/arch-course-labs/saga/pkg/payment/domain"
)

type PersistentProvider interface {
	PaymentRepository() domain.PaymentRepository
}

type UnitOfWork interface {
	Execute(f func(p PersistentProvider) error) error
}
