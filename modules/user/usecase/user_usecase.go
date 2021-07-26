package usecase

import (
	"context"
	"time"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/user"
	jwt "github.com/dgrijalva/jwt-go"
)

type userUsecase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

func NewUserUsecase(u user.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) fetchByEmail(c context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()

	res, err := u.userRepo.FetchByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userUsecase) FetchById(c context.Context, id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()

	res, err := u.userRepo.FetchById(ctx, id)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userUsecase) Register(c context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()

	existEmail, _ := u.fetchByEmail(ctx, user.Email)

	if existEmail != nil {
		return helpers.ErrConflict
	}

	user.Password = helpers.EncryptPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := u.userRepo.Store(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Login(c context.Context, user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()

	existEmail, err := u.userRepo.FetchByEmail(ctx, user.Email)

	if err != nil {
		return "", err
	}

	if existEmail == nil {
		return "", helpers.ErrConflict
	}

	user.Password = helpers.EncryptPassword(user.Password)

	if existEmail.Password != user.Password {
		return "", helpers.ErrWrongPassword
	}

	sign := jwt.New(jwt.SigningMethodHS256)

	claims := sign.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userid"] = existEmail.ID
	claims["email"] = existEmail.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// with hard-coded secret
	token, err := sign.SignedString([]byte("aqOeh4ck3R"))

	if err != nil {
		return "", err
	}

	return token, nil
}
