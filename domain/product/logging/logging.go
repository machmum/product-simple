package logging

import (
	"context"

	"github.com/elevenia/product-simple/domain/product"
	"github.com/elevenia/product-simple/internal/error"
	logger "github.com/elevenia/product-simple/internal/log"
	models "github.com/elevenia/product-simple/model"
)

type Product struct {
	uc product.Usecase
}

func NewProduct(usecase product.Usecase) *Product {
	return &Product{uc: usecase}
}

func (c *Product) View(ctx context.Context) (result *models.Product, err errin.Error) {
	defer func() {
		logger.ServiceLog(ctx, nil, nil, result, err)
	}()
	return c.uc.View(ctx)
}
