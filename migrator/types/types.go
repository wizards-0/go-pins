package types

const VERSION_SEPARATOR = "-"

type MigrationLog struct {
	Id int `db:"id" json:"id"`
	Migration
	Date int64  `db:"date" json:"date"`
	Hash string `db:"hash" json:"hash"`
}

type Migration struct {
	Name     string `db:"name" json:"name"`
	Version  string `db:"version" json:"version"`
	Query    string `db:"query" json:"query"`
	Rollback string `db:"rollback" json:"rollback"`
}
