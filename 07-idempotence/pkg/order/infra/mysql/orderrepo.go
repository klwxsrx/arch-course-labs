package mysql

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/infra/mysql"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/domain"
)

type repo struct {
	client mysql.Client
}

func (r *repo) NextID() uuid.UUID {
	return uuid.New()
}

func (r *repo) Store(order *domain.Order) error { // TODO: store order items
	const orderQuery = `
		INSERT INTO` + " `order` " + `(id, user_id, status, total_amount, created_at)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			user_id = user_id, status = status, total_amount = total_amount, updated_at = NOW()
	`

	binaryID, err := order.ID.MarshalBinary()
	if err != nil {
		return err
	}

	binaryUserID, err := order.UserID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = r.client.Exec(orderQuery, binaryID, binaryUserID, int(order.Status), order.TotalAmount)
	return err
}

func NewOrderRepository(client mysql.Client) domain.OrderRepository {
	return &repo{client: client}
}
