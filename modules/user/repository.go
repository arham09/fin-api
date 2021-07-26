package user

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Repository interface {
	FetchById(ctx context.Context, id int) (res *models.User, err error)
	FetchByEmail(ctx context.Context, email string) (res *models.User, err error)
	Store(ctx context.Context, u *models.User) error
}
