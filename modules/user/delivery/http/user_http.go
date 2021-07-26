package http

import (
	"context"
	"net/http"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/middleware"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHandler(e *echo.Echo, us user.Usecase, middleware *middleware.Middleware) {
	handler := &UserHandler{
		UserUsecase: us,
	}

	e.GET("/v1/profile", handler.FetchById, middleware.Authorize)
	e.POST("/v1/login", handler.Login)
	e.POST("/v1/register", handler.Register)
}

// ShowAccount godoc
// @Summary Show a user
// @Description get string by JWT token
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.User
// @Header 200 {string} Token "qwerty"
// @Router /profile [get]
func (u *UserHandler) FetchById(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get("userId").(int)

	if ctx == nil {
		ctx = context.Background()
	}

	user, err := u.UserUsecase.FetchById(ctx, userId)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// RegisterUser godoc
// @Summary Create a user
// @Description Create User
// @Accept  json
// @Produce  json
// @Param User body models.User true "models.User without ID"
// @Success 200 {object} models.User
// @Header 200 {string} Token "qwerty"
// @Router /register [post]
func (u *UserHandler) Register(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = helpers.VerifyEmail(user.Email); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	err = u.UserUsecase.Register(ctx, &user)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// LoginUser godoc
// @Summary Login user
// @Description Login User
// @Accept  json
// @Produce  json
// @Param User body models.User true "models.User without ID"
// @Success 200 {object} map[string]interface{}
// @Header 200 {string} Token "qwerty"
// @Router /login [post]
func (u *UserHandler) Login(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "email or password is missing",
		})
	}

	if err = helpers.VerifyEmail(user.Email); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	res, err := u.UserUsecase.Login(ctx, &user)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": res,
	})
}

func isRequestValid(m *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case helpers.ErrInternalServerError:
		return http.StatusInternalServerError
	case helpers.ErrNotFound:
		return http.StatusNotFound
	case helpers.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
