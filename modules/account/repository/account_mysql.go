package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/account"
	"github.com/sirupsen/logrus"
)

type mySqlAccountRepository struct {
	Conn *sql.DB
}

func NewMysqlAccountRepository(Conn *sql.DB) account.Repository {
	return &mySqlAccountRepository{Conn}
}

func (m *mySqlAccountRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Account, error) {
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

	result := make([]*models.Account, 0)

	for rows.Next() {
		t := new(models.Account)
		status := int(0)

		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Type,
			&t.Description,
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

func (m *mySqlAccountRepository) fetchTotal(ctx context.Context, query string) (int, error) {
	rows, err := m.Conn.QueryContext(ctx, query)

	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]int, 0)

	for rows.Next() {
		total := int(0)

		err = rows.Scan(
			&total,
		)

		if err != nil {
			logrus.Error(err)
			return 0, err
		}

		result = append(result, total)
	}

	return result[0], nil
}

func (m *mySqlAccountRepository) FetchAll(ctx context.Context, filters map[string]interface{}, keyword string, limit int, offset int) (res []*models.Account, total int, err error) {
	query := `SELECT id, name, type, description, status, created_at, updated_at FROM accounts WHERE status=1`
	countQuery := `SELECT COUNT(id) AS total FROM accounts WHERE status=1`

	for filter, param := range filters {
		addFilter := fmt.Sprintf(" AND %v = \"%v\"", filter, param)

		query = query + addFilter
		countQuery = countQuery + addFilter
	}

	if keyword != "" {
		addSearch := " AND name LIKE " + "'%" + keyword + "%'"
		query = query + addSearch
		countQuery = countQuery + addSearch
	}

	totalData, err := m.fetchTotal(ctx, countQuery)

	if err != nil {
		return nil, 0, err
	}

	pagination := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	query = query + pagination

	list, err := m.fetch(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	return list, totalData, nil
}

func (m *mySqlAccountRepository) FetchById(ctx context.Context, id int) (res *models.Account, err error) {
	query := `SELECT id, name, type, description, status, created_at, updated_at FROM accounts WHERE status=1 AND id = ?`

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

func (m *mySqlAccountRepository) Store(ctx context.Context, a *models.Account) error {
	query := `INSERT accounts SET name=?, type=?, description=?, created_at=?, updated_at=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.Type, a.Description, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	a.ID = int(lastID)

	return nil
}

func (m *mySqlAccountRepository) Update(ctx context.Context, a *models.Account) error {
	query := `UPDATE accounts SET name=?, type=?, description=?, updated_at=? WHERE status=1 AND id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.Type, a.Description, a.UpdatedAt, a.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

func (m *mySqlAccountRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE accounts SET status=0 WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
