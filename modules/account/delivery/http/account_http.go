package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/middleware"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/account"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type AccountHandler struct {
	AccountUsecase account.Usecase
}

func NewAccountHandler(e *echo.Echo, ac account.Usecase, middleware *middleware.Middleware) {
	handler := &AccountHandler{
		AccountUsecase: ac,
	}

	e.GET("/v1/account", handler.FetchAll, middleware.Authorize)
	e.GET("/v1/account/:id", handler.FetchById, middleware.Authorize)
	e.POST("/v1/account", handler.Create, middleware.Authorize)
	e.PATCH("/v1/account/:id", handler.Update, middleware.Authorize)
	e.DELETE("/v1/account/:id", handler.Delete, middleware.Authorize)
}

// ShowAccount godoc
// @Summary Show List account
// @Description get list account
// @Param keyword query string false "name search by keyword"
// @Param type query string false "filter by type"
// @Param limit query int true "limit list"
// @Param offset query int true "offset list"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Account in data
// @Header 200 {string} Token "qwerty"
// @Router /account [get]
func (a *AccountHandler) FetchAll(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	filter := make(map[string]interface{})
	keyword := ""

	params := c.QueryParams()

	for key, param := range params {
		if key != "limit" && key != "offset" {
			if key != "keyword" {
				filter[key] = param[0]
			} else {
				keyword = param[0]
			}
		}
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	res, total, err := a.AccountUsecase.FetchAll(ctx, filter, keyword, limit, offset)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  res,
		"total": total,
	})
}

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "account id"
// @Success 200 {object} models.Account
// @Header 200 {string} Token "qwerty"
// @Router /account/{id} [get]
func (a *AccountHandler) FetchById(c echo.Context) error {
	idAcc, err := strconv.Atoi(c.Param("id"))
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	user, err := a.AccountUsecase.FetchById(ctx, idAcc)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateAccount godoc
// @Summary Create an account
// @Description Create Account
// @Accept  json
// @Produce  json
// @Param account body models.Account true "models.Account without ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Account
// @Header 200 {string} Token "qwerty"
// @Router /account [post]
func (a *AccountHandler) Create(c echo.Context) error {
	var account models.Account

	err := c.Bind(&account)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&account); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AccountUsecase.Create(ctx, &account)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, account)
}

// DeleteAccount godoc
// @Summary Delete account
// @Description Delete account by ID
// @Accept  json
// @Produce  json
// @Param id path int true "account id"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200
// @Header 200 {string} Token "qwerty"
// @Router /account/{id} [delete]
func (a *AccountHandler) Delete(c echo.Context) error {
	idAcc, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ErrNotFound.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AccountUsecase.Delete(ctx, idAcc)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateAccount godoc
// @Summary Update account
// @Description Update account by ID
// @Accept  json
// @Produce  json
// @Param id path int true "account id"
// @Param account body models.Account true "models.Account without ID"
// @Success 200 {object} models.Account
// @Header 200 {string} Token "qwerty"
// @Router /account/{id} [patch]
func (a *AccountHandler) Update(c echo.Context) error {
	var account models.Account

	idAcc, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ErrNotFound.Error())
	}

	err = c.Bind(&account)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	account.ID = idAcc

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AccountUsecase.Update(ctx, &account)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func isRequestValid(m *models.Account) (bool, error) {
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
