package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/elevenia/product-simple/domain/product"
	repo "github.com/elevenia/product-simple/domain/product/repository"
	"github.com/elevenia/product-simple/helper/repository"
	"github.com/elevenia/product-simple/internal/error"
	"github.com/elevenia/product-simple/model"
)

type InternalDependency struct {
	// WithTrxFn represent wrapper func for transaction
	WithTrxFn repository.TransactionFn
}

func StdDependency() InternalDependency {
	return InternalDependency{WithTrxFn: repository.WithTransaction}
}

// NewProduct: will create new product object representation of product.Usecase interface
func NewProduct(conn *sql.DB, repo repo.Code, dep InternalDependency) product.Usecase {
	return &Product{
		pgDB:    conn,
		repo:    repo,
		timeout: time.Duration(1) * time.Second,
		dep:     dep,
	}
}

type Product struct {
	pgDB    *sql.DB
	repo    repo.Code
	timeout time.Duration
	mu      sync.Mutex
	dep     InternalDependency
}

func (c *Product) View(ctx context.Context) (*models.Product, errin.Error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var err error
	var result = new(models.Product)
	err = c.dep.WithTrxFn(c.pgDB, func(tx repository.Transaction) error {
		connTx := repository.Tx(tx)

		result, err = c.repo.View(ctxTimeout, connTx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errin.NewError(http.StatusInsufficientStorage, err)
	}

	return result, nil
}
