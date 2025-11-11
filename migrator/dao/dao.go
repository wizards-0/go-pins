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
	GetMigrationLogs(tx *sqlx.Tx) ([]types.MigrationLog, error)
	InsertMigrationLog(tx *sqlx.Tx, mLog types.MigrationLog) error
	DeleteMigrationLog(tx *sqlx.Tx, mLog types.MigrationLog) error
	ExecuteQuery(tx *sqlx.Tx, m types.Migration) error
	ExecuteRollback(tx *sqlx.Tx, m types.Migration) error
	SetupMigrationTable(tx *sqlx.Tx) error
}

type migrationDao struct {
	migrationTable string
}

func NewMigrationDao(schema string) MigrationDao {
	tableName := ""
	if schema != "" {
		tableName = schema + ".migration_log"
	} else {
		tableName = "migration_log"
	}

	return &migrationDao{
		migrationTable: tableName,
	}
}

func (dao *migrationDao) GetMigrationLogs(tx *sqlx.Tx) ([]types.MigrationLog, error) {
	mLogs := []types.MigrationLog{}

	if err := tx.Select(&mLogs, "SELECT id, name, version, query, rollback, date, hash FROM "+dao.migrationTable); err != nil {
		return nil, logger.WrapAndLogError(err, "error while getting migration logs from db")
	}

	sort.Slice(mLogs, func(i1, i2 int) bool {
		return semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version, types.VERSION_SEPARATOR)
	})
	return mLogs, nil
}

func (dao *migrationDao) InsertMigrationLog(tx *sqlx.Tx, mLog types.MigrationLog) error {
	_, err := tx.NamedExec("INSERT INTO "+dao.migrationTable+" (id, name, version, query, rollback, date, hash) VALUES (:id, :name, :version, :query, :rollback, :date, :hash)", &mLog)

	if err != nil {
		return logger.LogError(fmt.Errorf("error in database while inserting migration log\n%w", err))
	}
	return nil
}

func (dao *migrationDao) DeleteMigrationLog(tx *sqlx.Tx, mLog types.MigrationLog) error {
	_, err := tx.NamedExec("DELETE FROM "+dao.migrationTable+" WHERE version=:version", mLog)
	if err != nil {
		return logger.LogError(fmt.Errorf("error while deleting migration log\n%w", err))
	}
	return nil
}

func (dao *migrationDao) ExecuteQuery(tx *sqlx.Tx, m types.Migration) error {
	if _, err := tx.Exec(m.Query); err != nil {
		return logger.LogError(fmt.Errorf("error while executing query for migration '%v-%v'\n%w", m.Version, m.Name, err))
	}
	return nil
}

func (dao *migrationDao) ExecuteRollback(tx *sqlx.Tx, m types.Migration) error {
	if _, err := tx.Exec(m.Rollback); err != nil {
		return logger.LogError(fmt.Errorf("error while executing rollback query for version '%v'\n%w", m.Version, err))
	}
	return nil
}

func (dao *migrationDao) SetupMigrationTable(tx *sqlx.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS ` + dao.migrationTable + ` (
		id INTEGER PRIMARY KEY,
		name VARCHAR(200),
		version VARCHAR(20) UNIQUE,
		query TEXT,
		rollback TEXT,
		date BIGINT,
		hash VARCHAR(64)
	);`)
	if err != nil {
		return logger.LogError(fmt.Errorf("error in creating migration_log table\n%w", err))
	}
	return nil
}
