package persistence

import (
	"github.com/klwxsrx/arch-course-labs/saga/pkg/common/app/idempotence"
	"github.com/klwxsrx/arch-course-labs/saga/pkg/warehouse/domain"
)

type PersistentProvider interface {
	Stock() domain.Stock
	IdempotenceKeyStore() idempotence.KeyStore
}

type UnitOfWork interface {
	Execute(lockName string, f func(p PersistentProvider) error) error
}
