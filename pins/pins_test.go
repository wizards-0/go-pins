package pins

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
)

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	PanicOnError(nil)
	PanicOnError(errors.New(""))
}

func TestLog(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.Buffer{}
	logger.SetWriter(&buf, &buf, &buf, &buf)
	LogOnError(nil)
	assert.Equal("", buf.String())
	LogOnError(errors.New("test"))
	assert.Contains(buf.String(), "test")
	logger.ResetWriters()
}

func TestAssertValue(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	AssertValue("1", "1")
	AssertValue("1", "2")
}

func TestAppend(t *testing.T) {
	assert := assert.New(t)
	buf := &bytes.Buffer{}

	someStr := "someString"

	buf.Reset()
	sVal := "test value"
	AppendIfPresent(buf, sVal, someStr)
	assert.Equal("someString", buf.String())

	buf.Reset()
	sVal = ""
	AppendIfPresent(buf, sVal, someStr)
	assert.Equal("", buf.String())

	buf.Reset()
	iVal := 5
	AppendIfPresent(buf, iVal, someStr)
	assert.Equal("someString", buf.String())

	buf.Reset()
	iVal = 0
	AppendIfPresent(buf, iVal, someStr)
	assert.Equal("", buf.String())

	buf.Reset()
	bVal := false
	AppendIfPresent(buf, bVal, someStr)
	assert.Equal("someString", buf.String())

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

func TestMergeError(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(MergeErrors(nil))
	assert.Nil(MergeErrors())
	e1 := fmt.Errorf("msg1")
	e2 := fmt.Errorf("msg2")
	e3 := fmt.Errorf("msg3")
	assert.ErrorContains(MergeErrors(e1), "msg1")
	assert.ErrorContains(MergeErrors(e1, e2), "msg1. msg2")
	assert.ErrorContains(MergeErrors(e1, e2, e3), "msg1. msg2. msg3")
	assert.ErrorContains(MergeErrors(e1, nil), "msg1")
	assert.ErrorContains(MergeErrors(nil, e1, nil), "msg1")
	assert.ErrorContains(MergeErrors(e1, nil, e3), "msg1. msg3")
}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}
