package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elevenia/product-simple/config/env"
)

func InitDB(provider string) (dbConn *sql.DB) {
	var dsn string

	dbHost := cfg.Env.DBHost
	dbPort := cfg.Env.DBPort
	dbUser := cfg.Env.DBUser
	dbPass := cfg.Env.DBPass
	dbName := cfg.Env.DBName

	switch provider {
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	case "mysql":
		opt := "parseTime=1"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, opt)
	}
	fmt.Println(fmt.Sprintf("DB started on [%s:%d]@%s", dbHost, dbPort, dbName))

	dbConn, err := sql.Open(provider, dsn)
	if err != nil {
		fmt.Println(fmt.Sprintf("config: failed open DB: %v", err))
	}

	setConn(dbConn)

	err = dbConn.Ping()
	if err != nil {
		fmt.Println(fmt.Sprintf("config: failed ping DB: %v", err))
		os.Exit(1)
	}

	return
}

func setConn(dbConn *sql.DB) {
	var sb strings.Builder
	sb.WriteString("SetConnMaxLifetime=")

	// SetConnMaxLifetime
	dbConn.SetConnMaxLifetime(time.Duration(cfg.Env.DBConnMaxLifetime) * time.Minute)
	sb.WriteString(fmt.Sprintf("%dm", cfg.Env.DBConnMaxLifetime))

	sb.WriteString(". SetMaxOpenConns=")

	// SetMaxOpenConns.
	dbConn.SetMaxOpenConns(cfg.Env.DBConnMaxOpenConn)
	sb.WriteString(fmt.Sprintf("%d", cfg.Env.DBConnMaxOpenConn))

	sb.WriteString(". SetMaxIdleConns=")

	// SetMaxIdleConns.
	dbConn.SetMaxIdleConns(cfg.Env.DBConnMaxIdleConn)
	sb.WriteString(fmt.Sprintf("%d", cfg.Env.DBConnMaxIdleConn))

	fmt.Println(sb.String())
}
