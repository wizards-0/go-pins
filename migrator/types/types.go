package types

const VERSION_SEPARATOR = "-"

type MigrationStatus string

const (
	MIGRATION_STATUS_STARTED = "STARTED"
	MIGRATION_STATUS_SUCCESS = "SUCCESS"
	MIGRATION_STATUS_FAILED  = "FAILED"
)

type MigrationLog struct {
	Id int `db:"id" json:"id"`
	Migration
	Status string `db:"status" json:"status"`
	Date   int64  `db:"date" json:"date"`
	Hash   string `db:"hash" json:"hash"`
}

type Migration struct {
	Name     string `db:"name" json:"name"`
	Version  string `db:"version" json:"version"`
	Query    string `db:"query" json:"query"`
	Rollback string `db:"rollback" json:"rollback"`
}
