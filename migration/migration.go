package migration

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/semver"
)

func GetMigrationLogs(db *sqlx.DB) []MigrationLog {
	mLogs := []MigrationLog{}
	err := db.Select(&mLogs, "SELECT * FROM MIGRATION_LOG")
	checkError(err)
	sort.Slice(mLogs, func(i1, i2 int) bool {
		return semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version)
	})
	return mLogs
}

func Migrate(db *sqlx.DB, mArr []Migration) {
	setupMigrationTable(db)
	executeMigrationQueries(db, mArr)
}

func Rollback(db *sqlx.DB, ver string) {
	checkVersionExists(db, ver)
	mLogs := GetMigrationLogs(db)
	sort.Slice(mLogs, func(i1, i2 int) bool {
		return !semver.CompareSemver(mLogs[i1].Version, mLogs[i2].Version)
	})
	for _, mLog := range mLogs {
		db.MustExec(mLog.Rollback)
		deleteMigrationLog(db, &mLog)
		if semver.CompareSemver(mLog.Version, ver) {
			return
		}
	}
}

func setupMigrationTable(db *sqlx.DB) {
	db.MustExec(`CREATE TABLE IF NOT EXISTS MIGRATION_LOG (
		ID INTEGER PRIMARY KEY,
		NAME VARCHAR(200),
		VERSION VARCHAR(20) UNIQUE,
		QUERY TEXT,
		ROLLBACK TEXT,
		DATE INTEGER,
		HASH VARCHAR(64)
	);`)
}

func executeMigrationQueries(db *sqlx.DB, mArr []Migration) {
	mMap := getMigrationVersionMap(db)
	sort.Slice(mArr, func(i1, i2 int) bool {
		return semver.CompareSemver(mArr[i1].Version, mArr[i2].Version)
	})
	for _, m := range mArr {
		hash := hashQuery(m.Query)
		if mLog, exists := mMap[m.Version]; exists {
			validateHash(&mLog, hash)
		} else {
			db.MustExec(m.Query)
			insertMigrationLog(db, &m, hash)
		}
	}
}

func getMigrationVersionMap(db *sqlx.DB) map[string]MigrationLog {
	mLogs := GetMigrationLogs(db)
	mMap := map[string]MigrationLog{}
	for _, mLog := range mLogs {
		mMap[mLog.Version] = mLog
	}
	return mMap
}

func hashQuery(q string) string {
	hasher := sha256.New()
	hasher.Write([]byte(q))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func validateHash(m *MigrationLog, hash string) {
	if m.Hash != hash {
		msg := fmt.Sprintf(
			"DB Migration checksum failed for version %v,"+
				"please manually rollback the changes from this latest up to this version."+
				"And delete entries from migration_log table for the same", m.Version)
		logger.Error(msg)
		panic(msg)
	}
}

func insertMigrationLog(db *sqlx.DB, q *Migration, hash string) {
	migrationLog := MigrationLog{}
	migrationLog.Migration = *q
	migrationLog.Date = time.Now().UnixMilli()
	migrationLog.Hash = hash
	_, err := db.NamedExec("INSERT INTO MIGRATION_LOG (NAME, VERSION, QUERY, ROLLBACK, DATE, HASH) VALUES (:NAME, :VERSION, :QUERY, :ROLLBACK, :DATE, :HASH)", &migrationLog)
	checkError(err)
}

func deleteMigrationLog(db *sqlx.DB, mLog *MigrationLog) {
	_, err := db.NamedExec("DELETE FROM MIGRATION_LOG WHERE VERSION=:VERSION", mLog)
	checkError(err)
}

func checkVersionExists(db *sqlx.DB, ver string) {
	var i int
	rows := db.QueryRowx("SELECT 1 FROM MIGRATION_LOG WHERE VERSION=$1", ver)
	if err := rows.Scan(&i); err != nil {
		logger.Error(err)
		errMsg := "Unable to find version in migration log. Verify the version is correct and migration has been correctly setup"
		logger.Error(errMsg)
		panic(errMsg)
	}
}

func checkError(err any) {
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}

type MigrationLog struct {
	Id int `db:"ID" json:"id"`
	Migration
	Date int64  `db:"DATE" json:"date"`
	Hash string `db:"HASH" json:"hash"`
}

type Migration struct {
	Name     string `db:"NAME" json:"name"`
	Version  string `db:"VERSION" json:"version"`
	Query    string `db:"QUERY" json:"query"`
	Rollback string `db:"ROLLBACK" json:"rollback"`
}
