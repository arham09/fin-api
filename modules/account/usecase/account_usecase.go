package usecase

import (
	"context"
	"time"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/account"
)

type accountUsecase struct {
	accountRepo    account.Repository
	contextTimeout time.Duration
}

func NewAccountUsecase(a account.Repository, timeout time.Duration) account.Usecase {
	return &accountUsecase{
		accountRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *accountUsecase) FetchAll(c context.Context, filter map[string]interface{}, keyword string, limit int, offset int) ([]*models.Account, int, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	res, total, err := a.accountRepo.FetchAll(ctx, filter, keyword, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (a *accountUsecase) FetchById(c context.Context, id int) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	res, err := a.accountRepo.FetchById(ctx, id)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *accountUsecase) Delete(c context.Context, id int) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	existEmail, err := a.accountRepo.FetchById(ctx, id)

	if err != nil {
		return err
	}

	if existEmail == nil {
		return helpers.ErrNotFound
	}

	return a.accountRepo.Delete(ctx, id)
}

func (a *accountUsecase) Create(c context.Context, account *models.Account) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	err := a.accountRepo.Store(ctx, account)

	if err != nil {
		return err
	}
	return nil
}

func (a *accountUsecase) Update(c context.Context, account *models.Account) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	existID, err := a.accountRepo.FetchById(ctx, account.ID)

	if err != nil {
		return nil, err
	}

	if existID == nil {
		return nil, helpers.ErrNotFound
	}

	account.UpdatedAt = time.Now()

	err = a.accountRepo.Update(c, account)

	if err != nil {
		return nil, err
	}

	res, err := a.accountRepo.FetchById(c, account.ID)

	if err != nil {
		return nil, err
	}

	return res, nil
}
