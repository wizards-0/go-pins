package dao

import (
	"bytes"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/types"
)

var log = bytes.Buffer{}
var db *sqlx.DB
var dao MigrationDao

func setup() {
	db = getDbConnection()
	dao = NewMigrationDao(db, "")
	dao.SetupMigrationTable()
	logger.SetWriter(&log, &log, &log, &log)
}

func TestCrud(t *testing.T) {
	assert := assert.New(t)
	db = getDbConnection()
	defer db.Close()
	mDao := NewMigrationDao(db, "")
	mDao.SetupMigrationTable()

	var l1 = types.MigrationLog{Id: 1, Migration: types.Migration{Name: "Create test table", Version: "1"}}
	var l1o1 = types.MigrationLog{Id: 2, Migration: types.Migration{Name: "Add Column", Version: "1-1"}}
	var l2 = types.MigrationLog{Id: 3, Migration: types.Migration{Name: "Create test table2", Version: "2"}}

	_ = mDao.InsertMigrationLog(l1o1)
	_ = mDao.InsertMigrationLog(l2)
	_ = mDao.InsertMigrationLog(l1)
	mLogs, _ := mDao.GetMigrationLogs()
	assert.Equal(3, len(mLogs))
	assert.Equal("1", mLogs[0].Version)
	assert.Equal("1-1", mLogs[1].Version)
	assert.Equal("2", mLogs[2].Version)

	l1.Status = types.MIGRATION_STATUS_SUCCESS
	mDao.UpdateMigrationStatus(l1)
	mLogs, _ = mDao.GetMigrationLogs()
	assert.Equal(types.MIGRATION_STATUS_SUCCESS, mLogs[0].Status)

	mDao.DeleteMigrationLog(l1)
	mLogs, _ = mDao.GetMigrationLogs()
	assert.Equal(2, len(mLogs))

}

func TestSetupError(t *testing.T) {
	assert := assert.New(t)
	setup()
	db.Close()
	err := dao.SetupMigrationTable()
	assert.ErrorContains(err, "error in creating migration_log table")
}

func TestGetMigrationLogsError(t *testing.T) {
	assert := assert.New(t)
	setup()
	db.Close()
	_, err := dao.GetMigrationLogs()
	assert.ErrorContains(err, "error while getting migration logs")
}

func TestInsertMigrationLogsError(t *testing.T) {
	assert := assert.New(t)
	setup()
	db.Close()
	err := dao.InsertMigrationLog(types.MigrationLog{})
	assert.ErrorContains(err, "error in database while inserting migration log")
}

func TestUpdateMigrationStatusError(t *testing.T) {
	assert := assert.New(t)
	setup()
	db.Close()
	err := dao.UpdateMigrationStatus(types.MigrationLog{})
	assert.ErrorContains(err, "error while updating migration log status")
}

func TestDeleteMigrationLogError(t *testing.T) {
	assert := assert.New(t)
	setup()
	db.Close()
	err := dao.DeleteMigrationLog(types.MigrationLog{})
	assert.ErrorContains(err, "error while deleting migration log")
}

func TestExecQuery(t *testing.T) {
	assert := assert.New(t)
	setup()
	err := dao.ExecuteQuery(types.Migration{Query: "CREATE TABLE USER(ID INT,NAME TEXT)"})
	assert.Nil(err)
	_, err = db.Exec("SELECT * FROM USER")
	assert.Nil(err)
	db.Close()
	err = dao.ExecuteQuery(types.Migration{})
	assert.ErrorContains(err, "error while executing query")
}

func TestRollback(t *testing.T) {
	assert := assert.New(t)
	setup()
	err := dao.ExecuteRollback(types.Migration{Rollback: "DROP TABLE IF EXISTS USER"})
	assert.Nil(err)
	db.Close()
	err = dao.ExecuteRollback(types.Migration{})
	assert.ErrorContains(err, "error while executing rollback")
}

func TestSchema(t *testing.T) {
	assert := assert.New(t)
	dao := NewMigrationDao(nil, "").(*migrationDao)
	assert.Equal("migration_log", dao.migrationTable)

	dao = NewMigrationDao(nil, "my_schema").(*migrationDao)
	assert.Equal("my_schema.migration_log", dao.migrationTable)

}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}
