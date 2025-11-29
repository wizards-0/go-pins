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
	"github.com/samber/lo"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
	"github.com/wizards-0/go-pins/pins"
	"github.com/wizards-0/go-pins/semver"
)

type Migrator interface {
	Cli(osArgs []string) error
	GetMigrationLogs() ([]types.MigrationLog, error)
	RunMigrationsFromDirectory(path string) error
	Migrate(mArr []types.Migration) error
	Rollback(ver string) error
}

func New(db *sqlx.DB, schema string) Migrator {
	return &migrator{
		db:  db,
		dao: dao.NewMigrationDao(schema),
	}
}

func newMigrator(db *sqlx.DB, dao dao.MigrationDao) Migrator {
	return &migrator{
		db:  db,
		dao: dao,
	}
}

type migrator struct {
	db  *sqlx.DB
	dao dao.MigrationDao
}

func (m *migrator) Cli(osArgs []string) error {
	args := osArgs[1:]
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

func (m *migrator) GetMigrationLogs() (mArr []types.MigrationLog, err error) {
	txErr := pins.WithDefaultCtxTx(m.db, func(tx *sqlx.Tx) bool {
		mArr, err = m.dao.GetMigrationLogs(tx)
		return err == nil
	})
	return mArr, pins.MergeErrors(txErr, err)
}

func (m *migrator) RunMigrationsFromDirectory(path string) error {
	mArr, err := parseDirectory(path)
	if err != nil {
		return fmt.Errorf("error while running migrations from path %v\n%w", path, err)
	}
	return m.Migrate(mArr)
}

func (m *migrator) Migrate(mArr []types.Migration) error {
	var setupErr error
	txErr := pins.WithDefaultCtxTx(m.db, func(tx *sqlx.Tx) bool {
		setupErr = m.dao.SetupMigrationTable(tx)
		return setupErr == nil
	})
	err := pins.MergeErrors(txErr, setupErr)
	if err != nil {
		return fmt.Errorf("error while running migrations\n%w", err)
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
		if !semver.CompareSemver(ver, mLog.Version, types.VERSION_SEPARATOR) {
			return nil
		}

		var rollbackErr error
		txErr := pins.WithDefaultCtxTx(m.db, func(tx *sqlx.Tx) bool {
			if err := m.dao.ExecuteRollback(tx, mLog.Migration); err != nil {
				rollbackErr = fmt.Errorf("error while executing rollback query for version '%v'\n%w", ver, err)
				return false
			}

			if err := m.dao.DeleteMigrationLog(tx, mLog); err != nil {
				rollbackErr = fmt.Errorf("error while deleting migration log\n%w", err)
				return false
			}
			return true
		})
		err := pins.MergeErrors(txErr, rollbackErr)
		if err != nil {
			return err
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
	maxId := lo.MaxBy(lo.Values(mMap), func(mLog types.MigrationLog, maxLog types.MigrationLog) bool {
		return mLog.Id > maxLog.Id
	}).Id
	for _, m := range mArr {
		hash := hashQuery(m.Query)
		if mLog, exists := mMap[m.Version]; exists {
			if hashErr := validateHash(mLog, hash); hashErr != nil {
				return fmt.Errorf("error in execution while validating hash for '%v-%v'\n%w", mLog.Version, mLog.Name, hashErr)
			}
		} else {
			maxId = maxId + 1
			if execErr := migrator.executeQuery(m, maxId, hash); execErr != nil {
				return execErr
			}
		}
	}
	return nil
}

func (migrator migrator) executeQuery(m types.Migration, id int, hash string) error {
	var execErr error
	txErr := pins.WithDefaultCtxTx(migrator.db, func(tx *sqlx.Tx) bool {
		if err := migrator.dao.ExecuteQuery(tx, m); err != nil {
			execErr = logger.LogError(fmt.Errorf("error while executing query for migration '%v-%v'\n%w", m.Version, m.Name, err))
			return false
		}
		mLog, insertErr := migrator.insertMigrationLog(tx, m, id, hash)
		if insertErr != nil {
			execErr = logger.LogError(fmt.Errorf("error while inserting migration log for migration '%v-%v'\n%w", mLog.Version, mLog.Name, insertErr))
			return false
		}
		return true
	})

	return pins.MergeErrors(txErr, execErr)
}

func (m *migrator) getMigrationVersionMap() (mMap map[string]types.MigrationLog, err error) {
	txErr := pins.WithDefaultCtxTx(m.db, func(tx *sqlx.Tx) bool {
		mLogs, fetchErr := m.dao.GetMigrationLogs(tx)
		if fetchErr != nil {
			err = fmt.Errorf("error while getting version migrations map\n %w", fetchErr)
			return false
		}
		mMap = map[string]types.MigrationLog{}
		for _, mLog := range mLogs {
			mMap[mLog.Version] = mLog
		}
		return true
	})
	return mMap, pins.MergeErrors(txErr, err)
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

func (m *migrator) insertMigrationLog(tx *sqlx.Tx, q types.Migration, id int, hash string) (types.MigrationLog, error) {
	mLog := types.MigrationLog{}
	mLog.Id = id
	mLog.Migration = q
	mLog.Date = time.Now().UnixMilli()
	mLog.Hash = hash
	err := m.dao.InsertMigrationLog(tx, mLog)
	if err != nil {
		return types.MigrationLog{}, fmt.Errorf("error while inserting migration log\n%w", err)
	}
	return mLog, nil
}

func getMigrationInfo(mLogs []types.MigrationLog) string {
	buf := bytes.Buffer{}
	buf.WriteString("\n")
	for range 80 {
		buf.WriteRune('-')
	}
	buf.WriteString("\n")
	buf.WriteString("|  Version  ")
	writePadded(&buf, "|  Name", 80-13)
	buf.WriteString("|\n")
	for range 80 {
		buf.WriteRune('-')
	}
	buf.WriteString("\n")
	for _, m := range mLogs {
		writePadded(&buf, "|  "+m.Version, 12)
		writePadded(&buf, "|  "+m.Name, 80-13)
		buf.WriteString("|\n")
	}
	for range 80 {
		buf.WriteRune('-')
	}
	return buf.String()
}

func writePadded(buf *bytes.Buffer, s string, length int) {
	padLen := length - len(s)
	buf.WriteString(s)
	for range padLen {
		buf.WriteRune(' ')
	}
}
