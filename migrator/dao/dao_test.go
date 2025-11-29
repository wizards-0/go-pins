package dao

import (
	"bytes"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/types"
	"github.com/wizards-0/go-pins/pins"
)

var log = bytes.Buffer{}
var db *sqlx.DB
var dao MigrationDao

func setup() {
	db = getDbConnection()
	dao = NewMigrationDao("")
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		dao.SetupMigrationTable(tx)
		return true
	})

	logger.SetWriter(&log, &log, &log, &log)
}

func TestCrud(t *testing.T) {
	assert := assert.New(t)
	setup()

	var l1 = types.MigrationLog{Id: 1, Migration: types.Migration{Name: "Create test table", Version: "1"}}
	var l1o1 = types.MigrationLog{Id: 2, Migration: types.Migration{Name: "Add Column", Version: "1-1"}}
	var l2 = types.MigrationLog{Id: 3, Migration: types.Migration{Name: "Create test table2", Version: "2"}}

	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		_ = dao.InsertMigrationLog(tx, l1o1)
		_ = dao.InsertMigrationLog(tx, l2)
		_ = dao.InsertMigrationLog(tx, l1)
		mLogs, _ := dao.GetMigrationLogs(tx)
		assert.Equal(3, len(mLogs))
		assert.Equal("1", mLogs[0].Version)
		assert.Equal("1-1", mLogs[1].Version)
		assert.Equal("2", mLogs[2].Version)

		dao.DeleteMigrationLog(tx, l1)
		mLogs, _ = dao.GetMigrationLogs(tx)
		assert.Equal(2, len(mLogs))
		return false
	})

}

func TestSetupError(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		err := dao.SetupMigrationTable(tx)
		assert.ErrorContains(err, "error in creating migration_log table")
		return false
	})
}

func TestGetMigrationLogsError(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		_, err := dao.GetMigrationLogs(tx)
		assert.ErrorContains(err, "error while getting migration logs")
		return false
	})
}

func TestInsertMigrationLogsError(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		err := dao.InsertMigrationLog(tx, types.MigrationLog{})
		assert.ErrorContains(err, "error in database while inserting migration log")
		return false
	})
}

func TestDeleteMigrationLogError(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		tx.Rollback()
		err := dao.DeleteMigrationLog(tx, types.MigrationLog{})
		assert.ErrorContains(err, "error while deleting migration log")
		return false
	})
}

func TestExecQuery(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		err := dao.ExecuteQuery(tx, types.Migration{Query: "CREATE TABLE USER(ID INT,NAME TEXT)"})
		assert.Nil(err)
		_, err = tx.Exec("SELECT * FROM USER")
		assert.Nil(err)
		tx.Rollback()
		err = dao.ExecuteQuery(tx, types.Migration{})
		assert.ErrorContains(err, "error while executing query")
		return false
	})
}

func TestRollback(t *testing.T) {
	assert := assert.New(t)
	setup()
	pins.WithDefaultCtxTx(db, func(tx *sqlx.Tx) bool {
		err := dao.ExecuteRollback(tx, types.Migration{Rollback: "DROP TABLE IF EXISTS USER"})
		assert.Nil(err)
		tx.Rollback()
		err = dao.ExecuteRollback(tx, types.Migration{})
		assert.ErrorContains(err, "error while executing rollback")
		return false
	})
}

func TestSchema(t *testing.T) {
	assert := assert.New(t)
	dao := NewMigrationDao("").(*migrationDao)
	assert.Equal("migration_log", dao.migrationTable)

	dao = NewMigrationDao("my_schema").(*migrationDao)
	assert.Equal("my_schema.migration_log", dao.migrationTable)

}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}
