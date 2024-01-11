package config

import (
	"database/sql"
	"main/galats"
	"main/models"
	"time"

	sqlutil "github.com/AchmadRifai/sql-utils"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

func Str2Time(strTime string) (time.Time, error) {
	waktu, err := time.Parse(DDMMYYYYhhmmss, strTime)
	return waktu, err
}

func Time2Str(waktu time.Time) string {
	return waktu.Format(DDMMYYYYhhmmss)
}

func InitDb() {
	db, err := LoadDb()
	defer galats.DbClose(db)
	if err != nil {
		panic(err)
	}
	sqlutil.NewTransaction(db, true).Execute(func(db2 *sql.DB) {
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS type_account(name text,description text)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS account(name text,amount real,tipe text,increase real,decrease real,interest real,dueDate text,created text,description text)")
		sqlutil.DbSelect(db2, "CREATE TABLE IF NOT EXISTS transactions(id varchar primary key,description text,amount real,dari text,ke text,waktu text)")
		var types []models.TypeAccount
		types = append(types, models.TypeAccount{Name: "Account", Description: "Rekening"})
		types = append(types, models.TypeAccount{Name: "Asset", Description: "Aset"})
		types = append(types, models.TypeAccount{Name: "Income", Description: "Pemasukan"})
		types = append(types, models.TypeAccount{Name: "Liability", Description: "Kewajiban"})
		types = append(types, models.TypeAccount{Name: "Outcome", Description: "Pengeluaran"})
		for _, ta := range types {
			rows, _ := sqlutil.DbSelect(db2, "SELECT*FROM type_account WHERE name=?", ta.Name)
			if rows == nil {
				_, err := db2.Exec("INSERT INTO type_account VALUES(?,?)", ta.Name, ta.Description)
				if err != nil {
					panic(err)
				}
			}
		}
	})
}

func LoadDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared")
	return db, err
}
