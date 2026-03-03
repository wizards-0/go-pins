package slu

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func GetInsertPlaceholders(cols int, rows int) string {
	if cols < 1 || rows < 1 {
		return ""
	}
	placeholders := getPlaceholders(cols)
	values := strings.Builder{}
	values.WriteString(placeholders)
	for i := 1; i < rows; i++ {
		values.WriteString(",\n\t")
		values.WriteString(placeholders)
	}
	return values.String()
}

func GetDeletePlaceholders(rows int) string {
	return getPlaceholders(rows)
}

func getPlaceholders(n int) string {
	if n < 1 {
		return ""
	}
	placeholder := strings.Builder{}
	placeholder.WriteString("(?")
	for i := 1; i < n; i++ {
		placeholder.WriteString(",?")
	}
	placeholder.WriteString(")")
	return placeholder.String()
}

func FlattenStructs[T any](values []T, fieldAccessors []func(val T) any) []any {
	result := []any{}
	for _, v := range values {
		for _, accessor := range fieldAccessors {
			result = append(result, accessor(v))
		}
	}
	return result
}

func WithDefaultCtxTx(db *sqlx.DB, fn func(tx *sqlx.Tx) bool) error {
	return WithTx(context.Background(), db, fn)
}

func WithTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) bool) error {
	txOptions := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	return WithTxOptions(ctx, db, txOptions, fn)
}

func WithRoTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) bool) error {
	txOptions := &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}
	return WithTxOptions(ctx, db, txOptions, fn)
}

func WithTxOptions(ctx context.Context, db *sqlx.DB, txOptions *sql.TxOptions, fn func(tx *sqlx.Tx) bool) error {
	tx, err := db.BeginTxx(ctx, txOptions)
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
