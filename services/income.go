package services

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/galats"
	"main/models"
	"net/http"

	arrayutils "github.com/AchmadRifai/array-utils"
	mapsutils "github.com/AchmadRifai/maps-utils"
	sqlutils "github.com/AchmadRifai/sql-utils"
	"github.com/gorilla/mux"
)

func (i *IncomeRouter) Add(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	i.Lock()
	defer i.Done(w)
	var err error
	i.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	i.Tx, err = i.Db.Begin()
	if err != nil {
		panic(err)
	}
	var body models.Income
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	_, err = config.Str2Time(body.Created)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(i.Db, "SELECT*FROM account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Income exists"))
	}
	_, err = i.Db.Exec("INSERT INTO account(name,created,tipe) VALUES(?,?,?)", body.Name, body.Created, "income")
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (i *IncomeRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=?", "Income")
	if err := json.NewEncoder(w).Encode(&models.Incomes{
		Datas: arrayutils.Map(rows, func(v map[string]string, index int) models.Income {
			return models.Income{Name: v["name"], Created: v["created"]}
		}),
	}); err != nil {
		panic(err)
	}
}

func (i *IncomeRouter) Del(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	i.Lock()
	defer i.Done(w)
	var err error
	i.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	i.Tx, err = i.Db.Begin()
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Name not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(i.Db, "SELECT*FROM account WHERE name=? AND tipe=?", name, "Income")
	if rows == nil {
		panic(fmt.Errorf("Income %s not found", name))
	}
	_, err = i.Db.Exec("DELETE FROM account WHERE name=?", name)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func NewIncomeRouter() IncomeRouter {
	return IncomeRouter{}
}

type IncomeRouter struct {
	galats.PesimisticLocking
}
