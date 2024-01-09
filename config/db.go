package config

import (
	"database/sql"
	"main/galats"

	sqlutil "github.com/AchmadRifai/sql-utils"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

func InitDb() {
	db, err := LoadDb()
	defer galats.DbClose(db)
	if err != nil {
		panic(err)
	}
	sqlutil.NewTransaction(db, true).Execute(func(db2 *sql.DB) {
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS account(name text,amount real,description text)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS asset(name text,amount real,increase real,decrease real)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS liability(name text,amount real,interest real,dueDate text)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS outcome(name text,created text)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS transactions(id varchar primary key,description text,amount real,dari text,ke text,waktu text)")
	})
}

func LoadDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared")
	return db, err
}
