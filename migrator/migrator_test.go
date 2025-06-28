package migrator

import (
	"bytes"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
	mocks "github.com/wizards-0/go-pins/mocks/migrator/dao"
)

var TYPE_MIGRATION_LOG = mock.AnythingOfType("types.MigrationLog")
var TYPE_MIGRATION = mock.AnythingOfType("types.Migration")

var buf = bytes.Buffer{}
var db *sqlx.DB
var mDao dao.MigrationDao
var mRun Migrator

func setup() {
	w := io.Writer(&buf)
	logger.SetWriter(w, w, w, w)
	logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	db = getDbConnection()
	mDao = dao.NewMigrationDao(db)
	mRun = New(db)
}

func tearDown() {
	db.Close()
}

func TestMigration(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()

	mRun.Migrate([]types.Migration{q2, q1})
	mLogs, _ := mRun.GetMigrationLogs()
	assert.Equal(2, len(mLogs))
	assert.Equal("1", mLogs[0].Version)
	assert.Equal("2", mLogs[1].Version)

	mRun.Migrate([]types.Migration{q2, q1, q1_1})
	mLogs, _ = mRun.GetMigrationLogs()
	assert.Equal(3, len(mLogs))
	assert.Equal("1-1", mLogs[1].Version)
	assert.Equal(types.MIGRATION_STATUS_SUCCESS, mLogs[1].Status)
}

func TestMigrationFromDir(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()

	err := mRun.RunMigrationsFromDirectory("../resources/test/migrations/valid")
	assert.Nil(err)
	mLogs, _ := mRun.GetMigrationLogs()
	assert.Equal(1, len(mLogs))
	assert.Equal(types.MIGRATION_STATUS_SUCCESS, mLogs[0].Status)
}

func TestSetupDbError(t *testing.T) {
	setup()
	db.Close()
	err := mRun.Migrate([]types.Migration{q2, q1})
	assert.ErrorContains(t, err, "error in creating migration_log table")
}

func TestInvalidPathError(t *testing.T) {
	setup()
	defer tearDown()
	err := mRun.RunMigrationsFromDirectory("../non-existing-path")
	assert.ErrorContains(t, err, "error while running migrations from path")
}

func TestGetMigrationLogError(t *testing.T) {
	setup()
	db.Close()
	err := mRun.(*migrator).executeMigrationQueries([]types.Migration{q2, q1})
	assert.ErrorContains(t, err, "error while executing")
}

func TestInsertLogDbError(t *testing.T) {
	setup()
	defer tearDown()
	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)
	mockDao.EXPECT().InsertMigrationLog(mock.Anything).Return(-1, errors.New(""))
	err := mRun.(*migrator).executeQuery(q1, hashQuery(q1.Query))
	assert.ErrorContains(t, err, "error while inserting")
}

func TestInvalidHashDbError(t *testing.T) {
	setup()
	defer tearDown()
	_ = mRun.Migrate([]types.Migration{q1, q2})
	// Example depicting, if query is changed after being executed on db.
	// Causing hash to differ, throwing checksum error
	err := mRun.Migrate([]types.Migration{modifiedQ1, q2})
	assert.ErrorContains(t, err, "DB Migration checksum failed")
}

func TestRollback(t *testing.T) {
	setup()
	defer tearDown()
	mRun.Migrate([]types.Migration{q1, q2})

	mRun.Rollback("2")
	mLogs, _ := mRun.GetMigrationLogs()
	assert.Equal(t, 1, len(mLogs))
	assert.Equal(t, "1", mLogs[0].Version)

	mRun.Rollback("1")
	mLogs, _ = mRun.GetMigrationLogs()
	assert.Equal(t, 0, len(mLogs))
}

func TestRollbackWithoutMigration(t *testing.T) {
	setup()
	defer tearDown()
	mRun.Migrate([]types.Migration{})
	err := mRun.Rollback("2")
	assert.Nil(t, err)
}

func TestRollbackInvalidVersion(t *testing.T) {
	setup()
	defer tearDown()

	mRun.Migrate([]types.Migration{q1, q2})
	err := mRun.Rollback("1-1")

	assert.Nil(t, err)
	mLogs, _ := mRun.GetMigrationLogs()
	assert.Equal(t, 1, len(mLogs))
	assert.Equal(t, "1", mLogs[0].Version)

}

func TestExecQueryError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()

	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().SetupMigrationTable().RunAndReturn(func() error {
		mDao.SetupMigrationTable()
		return nil
	})
	mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
		return mDao.GetMigrationLogs()
	})
	mockDao.EXPECT().InsertMigrationLog(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) (int, error) {
		return mDao.InsertMigrationLog(mLog)
	})
	mockDao.EXPECT().UpdateMigrationStatus(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) error {
		return mDao.UpdateMigrationStatus(mLog)
	})
	mockDao.EXPECT().ExecuteQuery(TYPE_MIGRATION).Return(errors.New(""))

	mRun.Migrate([]types.Migration{q1})
	mLogs, _ := mRun.GetMigrationLogs()
	assert.Equal(types.MIGRATION_STATUS_FAILED, mLogs[0].Status)
}

func TestExecQuerySuccessStatusUpdateError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()

	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().SetupMigrationTable().RunAndReturn(func() error {
		mDao.SetupMigrationTable()
		return nil
	})
	mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
		return mDao.GetMigrationLogs()
	})
	mockDao.EXPECT().InsertMigrationLog(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) (int, error) {
		return mDao.InsertMigrationLog(mLog)
	})
	mockDao.EXPECT().UpdateMigrationStatus(TYPE_MIGRATION_LOG).Return(errors.New("1"))
	mockDao.EXPECT().ExecuteQuery(TYPE_MIGRATION).RunAndReturn(func(m types.Migration) error {
		return mDao.ExecuteQuery(m)
	})

	err := mRun.Migrate([]types.Migration{q1})
	assert.ErrorContains(err, "error while marking migration as success")
}

func TestExecQueryFailureStatusUpdateError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()

	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().SetupMigrationTable().RunAndReturn(func() error {
		mDao.SetupMigrationTable()
		return nil
	})
	mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
		return mDao.GetMigrationLogs()
	})
	mockDao.EXPECT().InsertMigrationLog(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) (int, error) {
		return mDao.InsertMigrationLog(mLog)
	})
	mockDao.EXPECT().UpdateMigrationStatus(TYPE_MIGRATION_LOG).Return(errors.New(""))
	mockDao.EXPECT().ExecuteQuery(TYPE_MIGRATION).Return(errors.New(""))

	err := mRun.Migrate([]types.Migration{q1})
	assert.ErrorContains(err, "error while marking migration as failed")
}

func TestRollbackFetchError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()
	mRun.Migrate([]types.Migration{q1})
	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().GetMigrationLogs().Return(nil, errors.New("fetch error"))

	err := mRun.Rollback("0")
	assert.ErrorContains(err, "fetch error")
}

func TestRollbackExecError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()
	mRun.Migrate([]types.Migration{q1})
	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
		return mDao.GetMigrationLogs()
	})
	mockDao.EXPECT().ExecuteRollback(TYPE_MIGRATION).Return(errors.New("roll back error"))

	err := mRun.Rollback("0")
	assert.ErrorContains(err, "error while executing rollback")
}

func TestRollbackDeleteLogError(t *testing.T) {
	assert := assert.New(t)
	setup()
	defer tearDown()
	mRun.Migrate([]types.Migration{q1})
	mockDao := mocks.NewMockMigrationDao(t)
	mRun = newMigrator(mockDao)

	mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
		return mDao.GetMigrationLogs()
	})
	mockDao.EXPECT().ExecuteRollback(TYPE_MIGRATION).RunAndReturn(func(m types.Migration) error {
		return mDao.ExecuteRollback(m)
	})
	mockDao.EXPECT().DeleteMigrationLog(TYPE_MIGRATION_LOG).Return(errors.New("delete error"))

	err := mRun.Rollback("0")
	assert.ErrorContains(err, "error while deleting migration log")
}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		log.Panic(err)
	}
	return db
}

var q1 = types.Migration{Name: "Create test table", Version: "1",
	Query:    "CREATE TABLE IF NOT EXISTS TEST(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST;",
}

var modifiedQ1 = types.Migration{Name: "Create test table1", Version: "1",
	Query:    "CREATE TABLE IF NOT EXISTS TEST1(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST1;",
}

var q2 = types.Migration{Name: "Create test table2", Version: "2",
	Query:    "CREATE TABLE IF NOT EXISTS TEST2(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST2;",
}

var q1_1 = types.Migration{Name: "Add Column", Version: "1-1",
	Query:    "ALTER TABLE TEST ADD COLUMN DESCRIPTION VARCHAR(2000)",
	Rollback: "ALTER TABLE TEST DROP COLUMN DESCRIPTION",
}
