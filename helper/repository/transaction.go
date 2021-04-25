package repository

import (
	"context"
	"database/sql"
	"errors"
)

type Conn interface {
	Transaction
}

type conn struct {
	tx Transaction
}

func (c conn) Prepare(query string) (*sql.Stmt, error) {
	return c.tx.Prepare(query)
}

func (c conn) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.tx.PrepareContext(ctx, query)
}

func (c conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.tx.Exec(query, args...)
}

func (c conn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.tx.ExecContext(ctx, query, args...)
}

func (c conn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.tx.Query(query, args...)
}

func (c conn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.tx.QueryContext(ctx, query, args...)
}

func (c conn) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.tx.QueryRow(query, args...)
}

func (c conn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.tx.QueryRowContext(ctx, query, args...)
}

func Tx(tx Transaction) Conn {
	return conn{tx: tx}
}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TransactionFn func(*sql.DB, func(Transaction) error) error

func WithTransaction(db *sql.DB, trans func(Transaction) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var errTx error
	defer func() {
		if p := recover(); p != nil {
			errTx = tx.Rollback()
			if errTx != nil {
				err = errTx
			}
		} else if err != nil {
			errTx = tx.Rollback()
			if errTx != nil {
				err = errTx
			}
		} else {
			errTx = tx.Commit()
			if errTx != nil {
				err = errTx
			}
		}
	}()

	err = trans(tx)

	return err
}

func ExecWithContext(ctx context.Context, conn Conn, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

//goland:noinspection GoErrorStringFormat
func ExecAffectedWithContext(ctx context.Context, conn Conn, query string, args ...interface{}) error {
	res, err := ExecWithContext(ctx, conn, query, args...)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		return errors.New("Unexpected error")
	}

	return nil
}
