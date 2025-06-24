package dao

import (
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/types"
	"github.com/wizards-0/go-pins/semver"
)

type MigrationDao interface {
	GetMigrationLogs() ([]types.MigrationLog, error)
	InsertMigrationLog(mLog types.MigrationLog) (int, error)
	UpdateMigrationStatus(mLog types.MigrationLog) error
	DeleteMigrationLog(mLog types.MigrationLog) error
	ExecuteQuery(m types.Migration) error
	ExecuteRollback(m types.Migration) error
	SetupMigrationTable() error
}

type migrationDao struct {
	db *sqlx.DB
}

func NewMigrationDao(db *sqlx.DB) MigrationDao {
	return migrationDao{
		db: db,
	}
}

func (dao migrationDao) GetMigrationLogs() ([]types.MigrationLog, error) {
	mLogs := []types.MigrationLog{}
	if err := dao.db.Select(&mLogs, "SELECT * FROM MIGRATION_LOG"); err != nil {
		return nil, logger.WrapAndLogError(err, "error while getting migration logs from db")
	}

	sort.Slice(mLogs, func(i1, i2 int) bool {
		return semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version, types.VERSION_SEPARATOR)
	})
	return mLogs, nil
}

func (dao migrationDao) InsertMigrationLog(mLog types.MigrationLog) (int, error) {
	_, err := dao.db.NamedExec("INSERT INTO MIGRATION_LOG (NAME, VERSION, QUERY, ROLLBACK, STATUS, DATE, HASH) VALUES (:NAME, :VERSION, :QUERY, :ROLLBACK, :STATUS, :DATE, :HASH)", &mLog)

	if err != nil {
		return -1, logger.LogError(fmt.Errorf("error in database while inserting migration log\n%w", err))
	}

	return dao.getMigrationLogId(mLog)
}

func (dao migrationDao) getMigrationLogId(mLog types.MigrationLog) (int, error) {
	id := -1
	err := dao.db.Get(&id, "SELECT ID FROM MIGRATION_LOG WHERE VERSION=$1", mLog.Version)
	if err != nil {
		return -1, logger.LogError(fmt.Errorf("error in database while getting migration log id\n%w", err))
	}
	return id, nil
}

func (dao migrationDao) UpdateMigrationStatus(mLog types.MigrationLog) error {
	_, err := dao.db.NamedExec("UPDATE MIGRATION_LOG SET STATUS=:STATUS WHERE ID=:ID", &mLog)
	if err != nil {
		return logger.LogError(fmt.Errorf("error while updating migration log status\n%w", err))
	}
	return nil
}

func (dao migrationDao) DeleteMigrationLog(mLog types.MigrationLog) error {
	_, err := dao.db.NamedExec("DELETE FROM MIGRATION_LOG WHERE VERSION=:VERSION", mLog)
	if err != nil {
		return logger.LogError(fmt.Errorf("error while deleting migration log\n%w", err))
	}
	return nil
}

func (dao migrationDao) ExecuteQuery(m types.Migration) error {
	if _, err := dao.db.Exec(m.Query); err != nil {
		return logger.LogError(fmt.Errorf("error while executing query for migration '%v-%v'\n%w", m.Version, m.Name, err))
	}
	return nil
}

func (dao migrationDao) ExecuteRollback(m types.Migration) error {
	if _, err := dao.db.Exec(m.Rollback); err != nil {
		return logger.LogError(fmt.Errorf("error while executing rollback query for version '%v'\n%w", m.Version, err))
	}
	return nil
}

func (dao migrationDao) SetupMigrationTable() error {
	_, err := dao.db.Exec(`CREATE TABLE IF NOT EXISTS MIGRATION_LOG (
		ID INTEGER PRIMARY KEY,
		NAME VARCHAR(200),
		VERSION VARCHAR(20) UNIQUE,
		QUERY TEXT,
		ROLLBACK TEXT,
		STATUS VARCHAR(20),
		DATE INTEGER,
		HASH VARCHAR(64)
	);`)
	if err != nil {
		return logger.LogError(fmt.Errorf("error in creating migration_log table\n%w", err))
	}
	return nil
}
