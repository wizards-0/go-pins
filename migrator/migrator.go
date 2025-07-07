package migrator

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
	"github.com/wizards-0/go-pins/semver"
)

type Migrator interface {
	ParseCmdArgs(args []string) error
	GetMigrationLogs() ([]types.MigrationLog, error)
	RunMigrationsFromDirectory(path string) error
	Migrate(mArr []types.Migration) error
	Rollback(ver string) error
}

func New(db *sqlx.DB) Migrator {
	return &migrator{
		dao: dao.NewMigrationDao(db),
	}
}

func newMigrator(dao dao.MigrationDao) Migrator {
	return &migrator{
		dao: dao,
	}
}

type migrator struct {
	dao dao.MigrationDao
}

func (m *migrator) ParseCmdArgs(args []string) error {
	cmd := args[0]
	switch cmd {
	case "run":
		return m.parseMigrationArgs(args)
	case "rollback":
		return m.parseRollbackArgs(args)
	default:
		return errors.New("invalid migration command. Valid options are 'run <path>' | 'rollback <version>'")
	}
}

func (m *migrator) parseRollbackArgs(args []string) error {
	if len(args) != 2 {
		return errors.New("rollback command needs to have version as second arg. Example 'rollback 1.1'")
	}
	version := args[1]
	if err := m.Rollback(version); err != nil {
		return err
	}
	mLogs, fetchErr := m.GetMigrationLogs()
	if fetchErr != nil {
		return logger.WrapAndLogError(fetchErr, "rollback completed, but error in fetching migration log")
	}
	logger.Info("Migration rollback completed. Following are the remaining migrations")
	logger.Info(getMigrationInfo(mLogs))
	return nil
}

func (m *migrator) parseMigrationArgs(args []string) error {
	if len(args) != 2 {
		return errors.New("migration run command needs to have path as second arg. Example 'run 1.1'")
	}
	path := args[1]
	if err := m.RunMigrationsFromDirectory(path); err != nil {
		return err
	}

	mLogs, fetchErr := m.GetMigrationLogs()
	if fetchErr != nil {
		return logger.WrapAndLogError(fetchErr, "migration completed, but error in fetching migration log")
	}

	logger.Info("Migration completed. Following are the migrations executed / verified")
	logger.Info(getMigrationInfo(mLogs))
	return nil
}

func (m *migrator) GetMigrationLogs() ([]types.MigrationLog, error) {
	return m.dao.GetMigrationLogs()
}

func (m *migrator) RunMigrationsFromDirectory(path string) error {
	mArr, err := parseDirectory(path)
	if err != nil {
		return fmt.Errorf("error while running migrations from path %v\n%w", path, err)
	}
	return m.Migrate(mArr)
}

func (m *migrator) Migrate(mArr []types.Migration) error {
	if setupErr := m.dao.SetupMigrationTable(); setupErr != nil {
		return fmt.Errorf("error while running migrations\n%w", setupErr)
	}
	return m.executeMigrationQueries(mArr)
}

func (m *migrator) Rollback(ver string) error {
	mLogs, fetchErr := m.GetMigrationLogs()
	if fetchErr != nil {
		return logger.WrapAndLogError(fetchErr, "error in executing rollback")
	}
	sort.Slice(mLogs, func(i1, i2 int) bool {
		return !semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version, types.VERSION_SEPARATOR)
	})
	for _, mLog := range mLogs {
		if err := m.dao.ExecuteRollback(mLog.Migration); err != nil {
			return fmt.Errorf("error while executing rollback query for version '%v'\n%w", ver, err)
		}

		if err := m.dao.DeleteMigrationLog(mLog); err != nil {
			return fmt.Errorf("error while deleting migration log\n%w", err)
		}
		if semver.CompareSemver(ver, mLog.Version, types.VERSION_SEPARATOR) {
			return nil
		}
	}
	return nil
}

func (migrator migrator) executeMigrationQueries(mArr []types.Migration) error {
	mMap, fetchErr := migrator.getMigrationVersionMap()
	if fetchErr != nil {
		return fmt.Errorf("error while executing migration queries\n%w", fetchErr)
	}
	sort.Slice(mArr, func(i1, i2 int) bool {
		return semver.CompareSemver(mArr[i1].Version, mArr[i2].Version, types.VERSION_SEPARATOR)
	})
	for _, m := range mArr {
		hash := hashQuery(m.Query)
		if mLog, exists := mMap[m.Version]; exists {
			if hashErr := validateHash(mLog, hash); hashErr != nil {
				return fmt.Errorf("error in execution while validating hash for '%v-%v'\n%w", mLog.Version, mLog.Name, hashErr)
			}
		} else if execErr := migrator.executeQuery(m, hash); execErr != nil {
			return execErr
		}
	}
	return nil
}

func (migrator migrator) executeQuery(m types.Migration, hash string) error {
	mLog, insertErr := migrator.insertMigrationLog(m, hash)
	if insertErr != nil {
		return logger.LogError(fmt.Errorf("error while inserting migration log for migration '%v-%v'\n%w", mLog.Version, mLog.Name, insertErr))
	}

	if err := migrator.dao.ExecuteQuery(m); err != nil {
		mLog.Status = types.MIGRATION_STATUS_FAILED
		if updateErr := migrator.dao.UpdateMigrationStatus(mLog); updateErr != nil {
			updateErr = fmt.Errorf("error while marking migration as failed '%v-%v'\n%w", mLog.Version, mLog.Name, updateErr)
			return logger.LogError(fmt.Errorf(updateErr.Error()+"\n%w", err))
		}
		return logger.LogError(fmt.Errorf("error while executing query for migration '%v-%v'\n%w", mLog.Version, mLog.Name, err))
	} else {
		mLog.Status = types.MIGRATION_STATUS_SUCCESS
		if updateErr := migrator.dao.UpdateMigrationStatus(mLog); updateErr != nil {
			updateErr = fmt.Errorf("error while marking migration as success '%v-%v'\n%w", mLog.Version, mLog.Name, updateErr)
			return logger.LogError(fmt.Errorf(updateErr.Error()+"\n%w", err))
		}
	}
	return nil
}

func (m *migrator) getMigrationVersionMap() (map[string]types.MigrationLog, error) {
	mLogs, err := m.dao.GetMigrationLogs()
	if err != nil {
		return nil, fmt.Errorf("error while getting version migrations map\n %w", err)
	}
	mMap := map[string]types.MigrationLog{}
	for _, mLog := range mLogs {
		mMap[mLog.Version] = mLog
	}
	return mMap, nil
}

func hashQuery(q string) string {
	hasher := sha256.New()
	hasher.Write([]byte(q))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func validateHash(m types.MigrationLog, hash string) error {
	if m.Hash != hash {
		return fmt.Errorf(
			"DB Migration checksum failed for version %v,"+
				"please manually rollback the changes from this latest up to this version."+
				"And delete entries from migration_log table for the same", m.Version)
	}
	return nil
}

func (m *migrator) insertMigrationLog(q types.Migration, hash string) (types.MigrationLog, error) {
	mLog := types.MigrationLog{}
	mLog.Migration = q
	mLog.Status = types.MIGRATION_STATUS_STARTED
	mLog.Date = time.Now().UnixMilli()
	mLog.Hash = hash
	id, err := m.dao.InsertMigrationLog(mLog)
	if err != nil {
		return types.MigrationLog{}, fmt.Errorf("error while inserting migration log\n%w", err)
	}
	mLog.Id = id
	return mLog, nil
}

func getMigrationInfo(mLogs []types.MigrationLog) string {
	buf := bytes.Buffer{}
	buf.WriteString("\n")
	for _, m := range mLogs {
		buf.WriteString("Name: " + m.Name + " | Version: " + m.Version)
	}
	return buf.String()
}
