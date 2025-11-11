package pins

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wizards-0/go-pins/logger"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogOnError(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func AssertValue(expected any, actual any) {
	if expected != actual {
		panic(fmt.Errorf("value mismatch, expected: %s, actual: %s", expected, actual))
	}
}

func AppendIfPresent(base *bytes.Buffer, val any, s string) {
	if sVal, ok := val.(string); ok && sVal == "" {
		return
	}
	if iVal, ok := val.(int); ok && iVal == 0 {
		return
	}
	base.Write([]byte(s))
}

func WithTx(db *sqlx.DB, fn func(tx *sqlx.Tx) bool) error {
	return WithCtxTx(context.Background(), db, fn)
}

func WithCtxTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) bool) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("error in starting transaction. %w", err)
	}
	if fn(tx) {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error in committing transaction. %w", err)
		}
	} else {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error in rolling back transaction. %w", err)
		}
	}
	return nil
}

func MergeErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	var m error = nil
	for _, err := range errs {
		if m == nil {
			m = err
		} else {
			if err != nil {
				m = fmt.Errorf("%w. %w", m, err)
			}
		}
	}
	return m
}
