package repo

import (
	"context"

	"github.com/elevenia/product-simple/helper/repository"
	"github.com/elevenia/product-simple/model"
)

type Code interface {
	View(ctx context.Context, conn repository.Conn) (*models.Product, error)
}
