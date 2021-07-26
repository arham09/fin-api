package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arham09/fin-api/helpers"
	"github.com/arham09/fin-api/middleware"
	"github.com/arham09/fin-api/models"
	"github.com/arham09/fin-api/modules/transaction"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type TrxRequest struct {
	Name        string  `json:"name" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	AccountID   int     `json:"accountId" validate:"required"`
}

type TrxHandler struct {
	TrxUsecase transaction.Usecase
}

func NewAccountHandler(e *echo.Echo, tr transaction.Usecase, middleware *middleware.Middleware) {
	handler := &TrxHandler{
		TrxUsecase: tr,
	}

	e.GET("/v1/transaction", handler.FetchAll, middleware.Authorize)
	e.GET("/v1/transaction/:id", handler.FetchById, middleware.Authorize)
	e.GET("/v1/transaction/daily", handler.FetchDailySummary, middleware.Authorize)
	e.GET("/v1/transaction/monthly", handler.FetchMonthlySummary, middleware.Authorize)
	e.POST("/v1/transaction", handler.Create, middleware.Authorize)
	e.PATCH("/v1/transaction/:id", handler.Update, middleware.Authorize)
	e.DELETE("/v1/transaction/:id", handler.Delete, middleware.Authorize)
}

// ShowTransaction godoc
// @Summary Show List Transaction
// @Description get list Transaction
// @Param keyword query string false "name search by keyword"
// @Param type query string false "filter by type"
// @Param accountId query int64 false "filter by type"
// @Param limit query int true "limit list"
// @Param offset query int true "offset list"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Transaction in data
// @Header 200 {string} Token "qwerty"
// @Router /transaction [get]
func (t *TrxHandler) FetchAll(c echo.Context) error {
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

	res, total, err := t.TrxUsecase.FetchAll(ctx, filter, keyword, limit, offset)

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

// ShowTransaction godoc
// @Summary Show a Transaction
// @Description get string by ID
// @ID get-trx-by-int
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "Transaction id"
// @Success 200 {object} models.Transaction
// @Header 200 {string} Token "qwerty"
// @Router /transaction/{id} [get]
func (t *TrxHandler) FetchById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	user, err := t.TrxUsecase.FetchById(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// ShowTransactionSummary godoc
// @Summary Show a Transaction Daily Summary
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {array} models.SummaryDaily
// @Header 200 {string} Token "qwerty"
// @Router /transaction/daily [get]
func (t *TrxHandler) FetchDailySummary(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	data, err := t.TrxUsecase.DailySummary(ctx)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, data)
}

// ShowTransactionSummary godoc
// @Summary Show a Transaction Monthly Summary
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {array} models.SummaryMonthly
// @Header 200 {string} Token "qwerty"
// @Router /transaction/monthly [get]
func (t *TrxHandler) FetchMonthlySummary(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	data, err := t.TrxUsecase.MonnthlySummary(ctx)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, data)
}

// CreateTransaction godoc
// @Summary Create a Transaction
// @Description Create Transaction
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param account body TrxRequest true "TrxRequest Body"
// @Success 200 {object} models.Transaction
// @Header 200 {string} Token "qwerty"
// @Router /transaction [post]
func (t *TrxHandler) Create(c echo.Context) error {
	var req TrxRequest
	var trx models.Transaction

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	err := c.Bind(&req)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&req); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	trx.Name = req.Name
	trx.Type = req.Type
	trx.Description = req.Description
	trx.Account.ID = req.AccountID

	if req.Type == "out" {
		trx.AmountOut = req.Amount
	} else if req.Type == "in" {
		trx.AmountIn = req.Amount
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Type should be in or out",
		})
	}

	err = t.TrxUsecase.Create(ctx, &trx)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Created",
	})
}

// DeleteTransaction godoc
// @Summary Delete Transaction
// @Description Delete Transaction by ID
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "Transaction id"
// @Success 200
// @Header 200 {string} Token "qwerty"
// @Router /transaction/{id} [delete]
func (t *TrxHandler) Delete(c echo.Context) error {
	idAcc, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ErrNotFound.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = t.TrxUsecase.Delete(ctx, idAcc)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateTransaction godoc
// @Summary Update Transaction
// @Description Update Transaction by ID
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "account id"
// @Param account body TrxRequest true "TrxRequest Body"
// @Success 200 {object} models.Transaction
// @Header 200 {string} Token "qwerty"
// @Router /transaction/{id} [patch]
func (t *TrxHandler) Update(c echo.Context) error {
	var req TrxRequest
	var trx models.Transaction

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	err = c.Bind(&req)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&req); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	trx.ID = id
	trx.Name = req.Name
	trx.Type = req.Type
	trx.Description = req.Description
	trx.Account.ID = req.AccountID
	trx.AmountOut = 0
	trx.AmountIn = 0

	if req.Type == "out" {
		trx.AmountOut = req.Amount
	} else if req.Type == "in" {
		trx.AmountIn = req.Amount
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Type should be in or out",
		})
	}

	res, err := t.TrxUsecase.Update(ctx, &trx)

	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func isRequestValid(m *TrxRequest) (bool, error) {
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
