package user

import (
	"context"

	"github.com/arham09/fin-api/models"
)

type Usecase interface {
	FetchById(ctx context.Context, id int) (res *models.User, err error)
	Register(c context.Context, user *models.User) error
	Login(c context.Context, user *models.User) (string, error)
}
