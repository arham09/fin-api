package transaction

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Repository interface {
	FetchAll(ctx context.Context, filters map[string]interface{}, keyword string, limit int, offset int) (res []*models.Transaction, total int, err error)
	FetchById(ctx context.Context, id int) (res *models.Transaction, err error)
	Store(ctx context.Context, t *models.Transaction) error
	Update(ctx context.Context, a *models.Transaction) error
	Delete(ctx context.Context, id int) error
	DailySummary(ctx context.Context) ([]*models.SummaryDaily, error)
	MonthlySummary(ctx context.Context) ([]*models.SummaryMonthly, error)
}
