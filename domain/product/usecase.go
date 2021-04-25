package product

import (
	"context"

	"github.com/elevenia/product-simple/internal/error"
	"github.com/elevenia/product-simple/model"
)

type Usecase interface {
	View(ctx context.Context) (*models.Product, errin.Error)
}
