package models

import (
	uuid "github.com/satori/go.uuid"

	"gopkg.in/guregu/null.v3"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"id"`
	SKU         string    `json:"sku" db:"sku"`
	Description string    `json:"description" db:"description"`
	CreatedBy   uuid.UUID `json:"-"`
	CreatedAt   null.Time `json:"-"`
}
