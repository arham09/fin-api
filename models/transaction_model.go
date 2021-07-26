package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	Description string    `json:"description" validate:"required"`
	AmountIn    float64   `json:"amountIn" validate:"required"`
	AmountOut   float64   `json:"amountOut" validate:"required"`
	Status      string    `json:"status"`
	Account     Account   `json:"account"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SummaryDaily struct {
	Day        int     `json:"day"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	AverageIn  float64 `json:"averageIn"`
	AverageOut float64 `json:"averageOut"`
}

type SummaryMonthly struct {
	Month      int     `json:"month"`
	AverageIn  float64 `json:"averageIn"`
	AverageOut float64 `json:"averageOut"`
	Year       int     `json:"year"`
}
