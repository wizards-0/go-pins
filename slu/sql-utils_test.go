package slu

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
)

func TestGetInsertPlaceholders(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("(?,?),\n\t(?,?)", GetInsertPlaceholders(2, 2))
	assert.Equal("", GetInsertPlaceholders(0, 2))
	assert.Equal("", GetInsertPlaceholders(2, 0))
	assert.Equal("(?),\n\t(?)", GetInsertPlaceholders(1, 2))
	assert.Equal("(?,?)", GetInsertPlaceholders(2, 1))
	assert.Equal("(?,?,?),\n\t(?,?,?),\n\t(?,?,?),\n\t(?,?,?)", GetInsertPlaceholders(3, 4))
}

func TestGetDeletePlaceholders(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("(?,?)", GetDeletePlaceholders(2))
	assert.Equal("", GetDeletePlaceholders(0))
	assert.Equal("(?)", GetDeletePlaceholders(1))
	assert.Equal("(?,?,?)", GetDeletePlaceholders(3))
}

func TestFlattenStructs(t *testing.T) {
	assert := assert.New(t)
	testStructs := []testStruct{
		{fieldA: "a1", fieldB: "b1"},
		{fieldA: "a2", fieldB: "b2"},
	}
	accessors := []func(val testStruct) any{
		func(val testStruct) any { return 1 },
		func(val testStruct) any { return val.fieldA },
		func(val testStruct) any { return val.fieldB },
	}
	result := FlattenStructs(testStructs, accessors)
	assert.Equal([]any{1, "a1", "b1", 1, "a2", "b2"}, result)
}

func TestWithTx(t *testing.T) {
	assert := assert.New(t)
	db := getDbConnection()
	defer db.Close()
	mDao := dao.NewMigrationDao("")
	l1 := types.MigrationLog{Id: 1, Migration: types.Migration{Name: "Create test table", Version: "1"}}
	l2 := types.MigrationLog{Id: 3, Migration: types.Migration{Name: "Create test table2", Version: "2"}}
	var mLogs []types.MigrationLog

	WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		mDao.SetupMigrationTable(tx)
		mDao.InsertMigrationLog(tx, l1)
		return true
	})
	WithRoTx(context.Background(), db, func(tx *sqlx.Tx) bool {
		mLogs, _ = mDao.GetMigrationLogs(tx)
		return true
	})
	assert.Equal(1, len(mLogs))

	WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		mDao.InsertMigrationLog(tx, l2)
		return false
	})
	WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		mLogs, _ = mDao.GetMigrationLogs(tx)
		return true
	})
	assert.Equal(1, len(mLogs))
}

func TestWithTxBeginError(t *testing.T) {
	assert := assert.New(t)
	db := getDbConnection()
	db.Close()

	txErr := WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		return true
	})
	assert.ErrorContains(txErr, "error in starting transaction")
}

func TestWithTxCommitError(t *testing.T) {
	assert := assert.New(t)
	db := getDbConnection()
	defer db.Close()
	txErr := WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		return true
	})
	assert.ErrorContains(txErr, "error in committing transaction")
}

func TestWithTxRollbackError(t *testing.T) {
	assert := assert.New(t)
	db := getDbConnection()
	defer db.Close()
	txErr := WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		return false
	})
	assert.ErrorContains(txErr, "error in rolling back transaction")
}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}

type testStruct struct {
	fieldA string
	fieldB string
}
