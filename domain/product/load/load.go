package product

import (
	"database/sql"

	"github.com/elevenia/product-simple/domain/product"
	"github.com/elevenia/product-simple/domain/product/logging"
	"github.com/elevenia/product-simple/domain/product/repository/pgsql"
	"github.com/elevenia/product-simple/domain/product/usecase"
)

func Load(dbConn *sql.DB) product.Usecase {
	repo := pgsql.NewProductRepo()
	service := logging.NewProduct(usecase.NewProduct(dbConn, repo, usecase.StdDependency()))

	return service
}
