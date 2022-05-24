package persistence

import "github.com/klwxsrx/arch-course-labs/saga/pkg/delivery/domain"

type PersistentProvider interface {
	DeliveryRepository() domain.DeliveryRepository
}

type UnitOfWork interface {
	Execute(f func(p PersistentProvider) error) error
}
