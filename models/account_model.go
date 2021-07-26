package models

import "time"

type Account struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
