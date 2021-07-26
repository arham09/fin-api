package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/transaction"
	"github.com/sirupsen/logrus"
)

type mySqlTrxRepository struct {
	Conn *sql.DB
}

func NewMysqlTrxRepository(Conn *sql.DB) transaction.Repository {
	return &mySqlTrxRepository{Conn}
}

func (m *mySqlTrxRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Transaction, error) {
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

	result := make([]*models.Transaction, 0)

	for rows.Next() {
		t := new(models.Transaction)
		status := int(0)
		accountStatus := int(0)

		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Type,
			&t.Description,
			&t.AmountIn,
			&t.AmountOut,
			&status,
			&t.Account.ID,
			&t.Account.Name,
			&t.Account.Description,
			&t.Account.Type,
			&accountStatus,
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

		if accountStatus == 1 {
			t.Account.Status = "active"
		} else {
			t.Account.Status = "inactive"
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mySqlTrxRepository) fetchTotal(ctx context.Context, query string) (int, error) {
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

func (m *mySqlTrxRepository) FetchById(ctx context.Context, id int) (res *models.Transaction, err error) {
	query := `SELECT t.id, t.name, t.type, t.description, t.amount_in, t.amount_out, t.status, t.account_id, a.name, a.type, a.description, a.status, t.created_at, t.updated_at FROM transactions t LEFT JOIN accounts a ON t.account_id=a.id WHERE t.status=1 AND t.id=?`

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

func (m *mySqlTrxRepository) FetchAll(ctx context.Context, filters map[string]interface{}, keyword string, limit int, offset int) (res []*models.Transaction, total int, err error) {
	query := `SELECT t.id, t.name, t.type, t.description, t.amount_in, t.amount_out, t.status, t.account_id, a.name, a.type, a.description, a.status, t.created_at, t.updated_at FROM transactions t LEFT JOIN accounts a ON t.account_id=a.id WHERE t.status=1`
	countQuery := `SELECT COUNT(t.id) AS total FROM transactions t LEFT JOIN accounts a ON t.account_id=a.id WHERE t.status=1`

	for filter, param := range filters {
		var addFilter string

		if filter == "accountId" {
			addFilter = fmt.Sprintf(" AND account_id = %v", param)
		} else {
			addFilter = fmt.Sprintf(" AND t.%v = \"%v\"", filter, param)
		}

		query = query + addFilter
		countQuery = countQuery + addFilter
	}

	if keyword != "" {
		addSearch := " AND t.name LIKE " + "'%" + keyword + "%'"
		query = query + addSearch
		countQuery = countQuery + addSearch
	}

	totalData, err := m.fetchTotal(ctx, countQuery)

	if err != nil {
		return nil, 0, err
	}

	list, err := m.fetch(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	return list, totalData, nil
}

func (m *mySqlTrxRepository) Store(ctx context.Context, t *models.Transaction) error {
	query := `INSERT transactions SET name=?, account_id=?, type=?, description=?, amount_in=?, amount_out=?, created_at=?, updated_at=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, t.Name, t.Account.ID, t.Type, t.Description, t.AmountIn, t.AmountOut, t.CreatedAt, t.UpdatedAt)

	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	t.ID = int(lastID)

	return nil
}

func (m *mySqlTrxRepository) Update(ctx context.Context, t *models.Transaction) error {
	query := `UPDATE transactions SET name=?, account_id=?, type=?, description=?, amount_in=?, amount_out=?, updated_at=? WHERE status=1 AND id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, t.Name, t.Account.ID, t.Type, t.Description, t.AmountIn, t.AmountOut, t.UpdatedAt, t.ID)
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

func (m *mySqlTrxRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE transactions SET status=0 WHERE id = ?`

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

func (m *mySqlTrxRepository) DailySummary(ctx context.Context) ([]*models.SummaryDaily, error) {
	query := `SELECT Avg(NULLIF(amount_in ,0)) as avgIn, Avg(NULLIF(amount_out ,0)) as avgOut, day(created_at) as day, month(created_at) as month, year(created_at) as year FROM transactions GROUP BY day(created_at), month(created_at), year(created_at);`
	rows, err := m.Conn.QueryContext(ctx, query)

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

	result := make([]*models.SummaryDaily, 0)

	for rows.Next() {
		t := new(models.SummaryDaily)

		err = rows.Scan(
			&t.AverageIn,
			&t.AverageOut,
			&t.Day,
			&t.Month,
			&t.Year,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mySqlTrxRepository) MonthlySummary(ctx context.Context) ([]*models.SummaryMonthly, error) {
	query := `SELECT Avg(NULLIF(amount_in ,0)) as avgIn, Avg(NULLIF(amount_out ,0)) as avgOut, month(created_at) as month, year(created_at) as year FROM transactions GROUP BY month(created_at), year(created_at);`
	rows, err := m.Conn.QueryContext(ctx, query)

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

	result := make([]*models.SummaryMonthly, 0)

	for rows.Next() {
		t := new(models.SummaryMonthly)

		err = rows.Scan(
			&t.AverageIn,
			&t.AverageOut,
			&t.Month,
			&t.Year,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}
