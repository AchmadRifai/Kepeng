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

func (o *OutcomeRouter) Add(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	o.Lock()
	defer o.Done(w)
	var err error
	o.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	o.Tx, err = o.Db.Begin()
	if err != nil {
		panic(err)
	}
	var body models.Outcome
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	if _, err = config.Str2Time(body.Created); err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(o.Db, "SELECT*FROM account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Outcome %s not found", body.Name))
	}
	if _, err = o.Db.Exec("INSERT INTO account(name,created,tipe) VALUES(?,?,?)", body.Name, body.Created, "Outcome"); err != nil {
		panic(err)
	}
	if err = json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (o *OutcomeRouter) Dell(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	o.Lock()
	defer o.Done(w)
	var err error
	o.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	o.Tx, err = o.Db.Begin()
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Name not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(o.Db, "SELECT*FROM account WHERE name=? AND tipe=?", name, "Outcome")
	if rows == nil {
		panic(fmt.Errorf("Outcome %s not found", name))
	}
	_, err = o.Db.Exec("DELETE FROM account WHERE name=? AND tipe=?", name, "Outcome")
	if err != nil {
		panic(err)
	}
	if err = json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (o *OutcomeRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=?", "Outcome")
	if err := json.NewEncoder(w).Encode(&models.Outcomes{
		Datas: arrayutils.Map(rows, func(v map[string]string, i int) models.Outcome {
			return models.Outcome{Name: v["name"], Created: v["created"]}
		}),
	}); err != nil {
		panic(err)
	}
}

func NewOutcome() OutcomeRouter {
	return OutcomeRouter{}
}

type OutcomeRouter struct {
	galats.PesimisticLocking
}
