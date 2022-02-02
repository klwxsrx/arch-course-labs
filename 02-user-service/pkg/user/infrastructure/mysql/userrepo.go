package mysql

import (
	"database/sql"
	"errors"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/app"
)

const mysqlDuplicateEntryErrorNumber = 1062

type userRepo struct {
	client mysql.Client
}

func (r *userRepo) Create(u *app.User) error {
	const query = "INSERT INTO user (id, email, first_name, last_name)" +
		"VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id = id"

	binaryID, err := u.ID.MarshalBinary()
	if err != nil {
		return err
	}
	result, err := r.client.Exec(query, binaryID, u.Email, u.FirstName, u.LastName)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return app.ErrUserAlreadyExists
	}
	return nil
}

func (r *userRepo) Get(id uuid.UUID) (*app.User, error) {
	const query = "SELECT * FROM user WHERE id = ?"

	binaryID, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var dbUser sqlxUser
	err = r.client.Get(&dbUser, query, binaryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &app.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}, nil
}

func (r *userRepo) Update(u *app.User) error {
	const query = "UPDATE user SET email = ?, first_name = ?, last_name = ? WHERE id = ?"

	binaryID, err := u.ID.MarshalBinary()
	if err != nil {
		return err
	}
	result, err := r.client.Exec(query, u.Email, u.FirstName, u.LastName, binaryID)
	if mysqlErr, ok := err.(*mysql2.MySQLError); ok && mysqlErr.Number == mysqlDuplicateEntryErrorNumber {
		return app.ErrUserWithSpecifiedEmailAlreadyExists
	}
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return app.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) Delete(id uuid.UUID) error {
	const query = "DELETE FROM user WHERE id = ?"

	binaryID, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	result, err := r.client.Exec(query, binaryID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return app.ErrUserNotFound
	}
	return nil
}

func NewUserRepository(client mysql.Client) app.UserRepository {
	return &userRepo{client}
}

type sqlxUser struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
}
