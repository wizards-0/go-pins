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
	InsertMigrationLog(mLog types.MigrationLog) error
	UpdateMigrationStatus(mLog types.MigrationLog) error
	DeleteMigrationLog(mLog types.MigrationLog) error
	ExecuteQuery(m types.Migration) error
	ExecuteRollback(m types.Migration) error
	SetupMigrationTable() error
}

type migrationDao struct {
	db             *sqlx.DB
	migrationTable string
}

func NewMigrationDao(db *sqlx.DB, schema string) MigrationDao {
	tableName := ""
	if schema != "" {
		tableName = schema + ".migration_log"
	} else {
		tableName = "migration_log"
	}

	return &migrationDao{
		db:             db,
		migrationTable: tableName,
	}
}

func (dao *migrationDao) GetMigrationLogs() ([]types.MigrationLog, error) {
	mLogs := []types.MigrationLog{}
	if err := dao.db.Select(&mLogs, "SELECT id, name, version, query, rollback, status, date, hash FROM "+dao.migrationTable); err != nil {
		return nil, logger.WrapAndLogError(err, "error while getting migration logs from db")
	}

	sort.Slice(mLogs, func(i1, i2 int) bool {
		return semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version, types.VERSION_SEPARATOR)
	})
	return mLogs, nil
}

func (dao *migrationDao) InsertMigrationLog(mLog types.MigrationLog) error {
	_, err := dao.db.NamedExec("INSERT INTO "+dao.migrationTable+" (id, name, version, query, rollback, status, date, hash) VALUES (:id, :name, :version, :query, :rollback, :status, :date, :hash)", &mLog)

	if err != nil {
		return logger.LogError(fmt.Errorf("error in database while inserting migration log\n%w", err))
	}
	return nil
}

func (dao *migrationDao) UpdateMigrationStatus(mLog types.MigrationLog) error {
	_, err := dao.db.NamedExec("UPDATE "+dao.migrationTable+" SET status=:status WHERE id=:id", &mLog)
	if err != nil {
		return logger.LogError(fmt.Errorf("error while updating migration log status\n%w", err))
	}
	return nil
}

func (dao *migrationDao) DeleteMigrationLog(mLog types.MigrationLog) error {
	_, err := dao.db.NamedExec("DELETE FROM "+dao.migrationTable+" WHERE version=:version", mLog)
	if err != nil {
		return logger.LogError(fmt.Errorf("error while deleting migration log\n%w", err))
	}
	return nil
}

func (dao *migrationDao) ExecuteQuery(m types.Migration) error {
	if _, err := dao.db.Exec(m.Query); err != nil {
		return logger.LogError(fmt.Errorf("error while executing query for migration '%v-%v'\n%w", m.Version, m.Name, err))
	}
	return nil
}

func (dao *migrationDao) ExecuteRollback(m types.Migration) error {
	if _, err := dao.db.Exec(m.Rollback); err != nil {
		return logger.LogError(fmt.Errorf("error while executing rollback query for version '%v'\n%w", m.Version, err))
	}
	return nil
}

func (dao *migrationDao) SetupMigrationTable() error {
	_, err := dao.db.Exec(`CREATE TABLE IF NOT EXISTS ` + dao.migrationTable + ` (
		id INTEGER PRIMARY KEY,
		name VARCHAR(200),
		version VARCHAR(20) UNIQUE,
		query TEXT,
		rollback TEXT,
		status VARCHAR(20),
		date BIGINT,
		hash VARCHAR(64)
	);`)
	if err != nil {
		return logger.LogError(fmt.Errorf("error in creating migration_log table\n%w", err))
	}
	return nil
}
