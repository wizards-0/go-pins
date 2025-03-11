package migration

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Migrate() (*[]MigrationLog, error) {
	db := getDbConnection()
	defer db.Close()
	setupMigrationTable(db)
	executeMigrationQueries(db)
	return getMigrationLog(db), nil
}

func getDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:albion-mat-db?mode=memory&cache=shared")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func setupMigrationTable(db *sqlx.DB) {
	db.MustExec(`CREATE TABLE IF NOT EXISTS MIGRATION_LOG (
		ID INTEGER PRIMARY KEY,
		NAME VARCHAR(200),
		VERSION VARCHAR(20) UNIQUE,
		DATE INTEGER,
		HASH VARCHAR(64)
	);`)
}

func executeMigrationQueries(db *sqlx.DB) {
	appliedMigrations := getMigrationVersionMap(db)
	sort.Slice(migrationQueries, func(i1, i2 int) bool {
		return migrationQueries[i1].Version > migrationQueries[i2].Version
	})
	for _, migrationQuery := range migrationQueries {
		hash := hashQuery(migrationQuery.Query)
		if appliedMigration, versionApplied := appliedMigrations[migrationQuery.Version]; versionApplied {
			validateHash(&appliedMigration, hash)
		} else {
			db.MustExec(migrationQuery.Query)
			insertMigrationLog(db, migrationQuery, hash)
		}
	}
}

func getMigrationLog(db *sqlx.DB) *[]MigrationLog {
	migrations := &[]MigrationLog{}
	if err := db.Select(migrations, "SELECT * FROM MIGRATION"); err != nil {
		log.Panic(err)
	}
	return migrations
}

func getMigrationVersionMap(db *sqlx.DB) map[string]MigrationLog {
	appliedMigrations := getMigrationLog(db)
	appliedMigrationVersionMap := map[string]MigrationLog{}
	for _, appliedMigration := range *appliedMigrations {
		appliedMigrationVersionMap[appliedMigration.Version] = appliedMigration
	}
	return appliedMigrationVersionMap
}

func hashQuery(q string) string {
	hasher := sha256.New()
	hasher.Write([]byte(q))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func validateHash(m *MigrationLog, hash string) {
	if m.Hash != hash {
		log.Panic("DB Migration checksum failed, please manually rollback the changes and clear migration_log table")
	}
}

func insertMigrationLog(db *sqlx.DB, q Migration, hash string) {
	migration := &MigrationLog{
		-1,
		q.Name,
		q.Version,
		time.Now().UnixMilli(),
		hash,
	}
	_, err := db.NamedExec("INSERT INTO MIGRATION_LOG (NAME, VERSION, DATE, HASH) VALUES (:NAME, :VERSION, :DATE, :HASH)", migration)
	if err != nil {
		log.Panic(err)
	}
}

type MigrationLog struct {
	Id      int    `db:"ID" json:"id"`
	Name    string `db:"NAME" json:"name"`
	Version string `db:"VERSION" json:"version"`
	Date    int64  `db:"DATE" json:"date"`
	Hash    string `db:"HASH" json:"hash"`
}

type Migration struct {
	Name     string
	Version  string
	Query    string
	Rollback string
}

var migrationQueries = []Migration{
	{Name: "Create test table", Version: "1",
		Query: "CREATE TABLE IF NOT EXISTS TEST(Id int);",
	},
}
