package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/user"
	"github.com/sirupsen/logrus"
)

type mySqlUserRepository struct {
	Conn *sql.DB
}

func NewMysqlUserRepository(Conn *sql.DB) user.Repository {
	return &mySqlUserRepository{Conn}
}

func (m *mySqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.User, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.User, 0)

	for rows.Next() {
		t := new(models.User)
		status := int(0)

		err = rows.Scan(
			&t.ID,
			&t.Email,
			&t.Password,
			&t.Name,
			&status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		if status == 1 {
			t.Status = "active"
		} else {
			t.Status = "inactive"
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mySqlUserRepository) FetchById(ctx context.Context, id int) (res *models.User, err error) {
	query := `SELECT id, email, password, name, status, created_at, updated_at FROM users WHERE id = ?`

	list, err := m.fetch(ctx, query, id)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, helpers.ErrNotFound
	}

	return res, nil
}

func (m *mySqlUserRepository) FetchByEmail(ctx context.Context, email string) (res *models.User, err error) {
	query := `SELECT id, email, password, name, status, created_at, updated_at FROM users WHERE email = ?`

	list, err := m.fetch(ctx, query, email)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, helpers.ErrNotFound
	}

	fmt.Println(res)

	return res, nil
}

func (m *mySqlUserRepository) Store(ctx context.Context, u *models.User) error {
	query := `INSERT users SET email=?, password=?, name=?, created_at=?, updated_at=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, u.Email, u.Password, u.Name, u.CreatedAt, u.UpdatedAt)

	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	u.ID = int(lastID)

	return nil
}
