package account

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Usecase interface {
	// FetchAll(ctx context.Context, search string, limit int, offset int) (res *models.Account, err error)
	FetchAll(c context.Context, filter map[string]interface{}, keyword string, limit int, offset int) ([]*models.Account, int, error)
	FetchById(ctx context.Context, id int) (res *models.Account, err error)
	Create(ctx context.Context, a *models.Account) error
	Update(ctx context.Context, a *models.Account) (*models.Account, error)
	Delete(ctx context.Context, id int) error
}
