package types

const VERSION_SEPARATOR = "-"

type MigrationStatus string

const (
	MIGRATION_STATUS_STARTED = "STARTED"
	MIGRATION_STATUS_SUCCESS = "SUCCESS"
	MIGRATION_STATUS_FAILED  = "FAILED"
)

type MigrationLog struct {
	Id int `db:"ID" json:"id"`
	Migration
	Status string `db:"STATUS" json:"status"`
	Date   int64  `db:"DATE" json:"date"`
	Hash   string `db:"HASH" json:"hash"`
}

type Migration struct {
	Name     string `db:"NAME" json:"name"`
	Version  string `db:"VERSION" json:"version"`
	Query    string `db:"QUERY" json:"query"`
	Rollback string `db:"ROLLBACK" json:"rollback"`
}
