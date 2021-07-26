package usecase

import (
	"context"
	"time"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/account"
	"github.com/arham09/fin-api/modules/transaction"
)

type transactionUsecase struct {
	trxRepo        transaction.Repository
	accountRepo    account.Repository
	contextTimeout time.Duration
}

func NewTrxRepo(t transaction.Repository, a account.Repository, timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		trxRepo:        t,
		accountRepo:    a,
		contextTimeout: timeout,
	}
}

func (t *transactionUsecase) FetchAll(c context.Context, filter map[string]interface{}, keyword string, limit int, offset int) ([]*models.Transaction, int, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	res, total, err := t.trxRepo.FetchAll(ctx, filter, keyword, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (t *transactionUsecase) FetchById(c context.Context, id int) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	res, err := t.trxRepo.FetchById(ctx, id)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionUsecase) Delete(c context.Context, id int) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	existEmail, err := t.trxRepo.FetchById(ctx, id)

	if err != nil {
		return err
	}

	if existEmail == nil {
		return helpers.ErrNotFound
	}

	return t.trxRepo.Delete(ctx, id)
}

func (t *transactionUsecase) Create(c context.Context, trx *models.Transaction) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	trx.CreatedAt = time.Now()
	trx.UpdatedAt = time.Now()

	err := t.trxRepo.Store(ctx, trx)

	if err != nil {
		return err
	}
	return nil
}

func (t *transactionUsecase) Update(c context.Context, trx *models.Transaction) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	existID, err := t.trxRepo.FetchById(ctx, trx.ID)

	if err != nil {
		return nil, err
	}

	if existID == nil {
		return nil, helpers.ErrNotFound
	}

	accountId, err := t.accountRepo.FetchById(ctx, trx.Account.ID)

	if err != nil {
		return nil, err
	}

	if accountId == nil {
		return nil, helpers.ErrNotFound
	}

	trx.UpdatedAt = time.Now()

	err = t.trxRepo.Update(c, trx)

	if err != nil {
		return nil, err
	}

	res, err := t.trxRepo.FetchById(c, trx.ID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionUsecase) DailySummary(c context.Context) ([]*models.SummaryDaily, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	res, err := t.trxRepo.DailySummary(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionUsecase) MonnthlySummary(c context.Context) ([]*models.SummaryMonthly, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)

	defer cancel()

	res, err := t.trxRepo.MonthlySummary(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}
