package account

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Repository interface {
	FetchAll(ctx context.Context, filters map[string]interface{}, keyword string, limit int, offset int) (res []*models.Account, total int, err error)
	FetchById(ctx context.Context, id int) (res *models.Account, err error)
	Store(ctx context.Context, a *models.Account) error
	Update(ctx context.Context, a *models.Account) error
	Delete(ctx context.Context, id int) error
}
