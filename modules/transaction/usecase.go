package transaction

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Usecase interface {
	FetchAll(c context.Context, filter map[string]interface{}, keyword string, limit int, offset int) ([]*models.Transaction, int, error)
	FetchById(c context.Context, id int) (*models.Transaction, error)
	Create(c context.Context, trx *models.Transaction) error
	Update(c context.Context, trx *models.Transaction) (*models.Transaction, error)
	Delete(c context.Context, id int) error
	DailySummary(c context.Context) ([]*models.SummaryDaily, error)
	MonnthlySummary(c context.Context) ([]*models.SummaryMonthly, error)
}
