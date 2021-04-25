package pgsql

import (
	"context"
	"strings"

	repo "github.com/elevenia/product-simple/domain/product/repository"
	"github.com/elevenia/product-simple/helper/repository"
	logger "github.com/elevenia/product-simple/internal/log"
	"github.com/elevenia/product-simple/model"
)

type productRepo struct{}

// NewProductRepo: will create an object that represent the product.Repository interface
func NewProductRepo() repo.Code {
	return &productRepo{}
}

func (r *productRepo) View(ctx context.Context, conn repository.Conn) (*models.Product, error) {
	var err error
	var args []interface{}
	var sb strings.Builder

	defer func() {
		if err != nil {
			logger.SqlLog(ctx, logger.QueryHelper{Error: err}, sb.String(), args...)
		}
	}()

	sb.WriteString("SELECT id, name, sku, description ")
	sb.WriteString("FROM ")
	sb.WriteString(models.TblProduct)
	sb.WriteString(" limit 1")

	result := new(models.Product)
	err = conn.QueryRowContext(ctx, sb.String(), args...).Scan(
		&result.ID,
		&result.Name,
		&result.SKU,
		&result.Description,
	)
	return result, err
}
