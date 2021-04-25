package config

import (
	"database/sql"

	"github.com/elevenia/product-simple/config/env"
	"github.com/elevenia/product-simple/domain/product/handler/http"
	product "github.com/elevenia/product-simple/domain/product/load"
	"github.com/elevenia/product-simple/helper/response"
	"github.com/elevenia/product-simple/internal/log"
	"github.com/elevenia/product-simple/internal/server"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Load() {
	// load env variable
	cfg.LoadEnv()

	dbConn := InitDB(cfg.Env.DBProvider)
	defer dbConn.Close()

	// load log
	// set log level at debug, if env.debug true
	lvl := logrus.DebugLevel
	if !cfg.Env.Debug {
		lvl = logrus.InfoLevel
	}
	logger.New(logger.Options{Level: lvl})

	srv := srvin.New(ServerOpts())

	handler(srv.Echo(), dbConn)

	srv.Start()
}

func handler(e *echo.Echo, dbConn *sql.DB) {
	// initialize response
	response.Init(cfg.Env.Debug)

	http.NewHandler(e, product.Load(dbConn))
}
