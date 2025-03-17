package migration

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
)

var buf = bytes.Buffer{}

func setup() {
	w := io.Writer(&buf)
	logger.SetWriter(w, w, w, w)
}

func TestMigration(t *testing.T) {
	assert := assert.New(t)
	db := getDbConnection()
	defer db.Close()

	Migrate(db, []Migration{q2, q1})
	mLogs := GetMigrationLogs(db)
	assert.Equal(2, len(mLogs))
	assert.Equal("1", mLogs[0].Version)
	assert.Equal("2", mLogs[1].Version)

	Migrate(db, []Migration{q2, q1, q1_1})
	mLogs = GetMigrationLogs(db)
	assert.Equal(3, len(mLogs))
	assert.Equal("1.1", mLogs[1].Version)
}

func TestSetupDbError(t *testing.T) {
	setup()
	db := getDbConnection()
	db.Close()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	setupMigrationTable(db)
}

func TestExecDbError(t *testing.T) {
	setup()
	db := getDbConnection()
	db.Close()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	executeMigrationQueries(db, []Migration{q2, q1})
}

func TestInsertLogDbError(t *testing.T) {
	setup()
	db := getDbConnection()
	db.Close()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	insertMigrationLog(db, &q1, "")
}

func TestInvalidHashDbError(t *testing.T) {
	setup()
	db := getDbConnection()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		db.Close()
	}()
	Migrate(db, []Migration{q1, q2})
	// Example depicting, if query is changed after being executed on db.
	// Causing hash to differ, throwing checksum error
	Migrate(db, []Migration{modifiedQ1, q2})
}

func TestRollback(t *testing.T) {
	db := getDbConnection()
	defer db.Close()
	Migrate(db, []Migration{q1, q2})

	Rollback(db, "2")
	mLogs := GetMigrationLogs(db)
	assert.Equal(t, 1, len(mLogs))
	assert.Equal(t, "1", mLogs[0].Version)

	Rollback(db, "1")
	mLogs = GetMigrationLogs(db)
	assert.Equal(t, 0, len(mLogs))
}

func TestRollbackWithoutMigration(t *testing.T) {
	setup()
	db := getDbConnection()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		db.Close()
	}()
	Rollback(db, "2")
}

func TestRollbackInvalidVersion(t *testing.T) {
	setup()
	db := getDbConnection()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		db.Close()
	}()
	Migrate(db, []Migration{q1, q2})
	Rollback(db, "1.1")
}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test-db?mode=memory&cache=shared")
	if err != nil {
		log.Panic(err)
	}
	return db
}

var q1 = Migration{Name: "Create test table", Version: "1",
	Query:    "CREATE TABLE IF NOT EXISTS TEST(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST;",
}

var modifiedQ1 = Migration{Name: "Create test table1", Version: "1",
	Query:    "CREATE TABLE IF NOT EXISTS TEST1(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST1;",
}

var q2 = Migration{Name: "Create test table2", Version: "2",
	Query:    "CREATE TABLE IF NOT EXISTS TEST2(Id int);",
	Rollback: "DROP TABLE IF EXISTS TEST2;",
}

var q1_1 = Migration{Name: "Add Column", Version: "1.1",
	Query:    "ALTER TABLE TEST ADD COLUMN DESCRIPTION VARCHAR(2000)",
	Rollback: "ALTER TABLE TEST DROP COLUMN DESCRIPTION",
}
