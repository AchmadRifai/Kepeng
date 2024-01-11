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

func (t *TypeRouter) Add(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	t.Lock()
	var err error
	t.Db, err = config.LoadDb()
	defer t.Done(w)
	if err != nil {
		panic(err)
	}
	t.Tx, err = t.Db.Begin()
	if err != nil {
		panic(err)
	}
	var body models.TypeAccount
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(t.Db, "SELECT*FROM type_account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Type %s exists", body.Name))
	}
	_, err = t.Db.Exec("INSERT INTO type_account VALUES(?,?)", body.Name, body.Description)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (t *TypeRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM type_account")
	if err := json.NewEncoder(w).Encode(&models.TypeAccounts{
		Datas: arrayutils.Map(rows, func(v map[string]string, index int) models.TypeAccount {
			return models.TypeAccount{Name: v["name"], Description: v["description"]}
		}),
	}); err != nil {
		panic(err)
	}
}

func (t *TypeRouter) Del(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	t.Lock()
	var err error
	t.Db, err = config.LoadDb()
	defer t.Done(w)
	if err != nil {
		panic(err)
	}
	t.Tx, err = t.Db.Begin()
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Param name not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(t.Db, "SELECT*FROM type_account WHERE name=?", name)
	if rows == nil {
		panic(fmt.Errorf("Type %s not found", name))
	}
	_, err = t.Db.Exec("DELETE FROM type_account WHERE name=?", name)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func NewTypeRouter() TypeRouter {
	return TypeRouter{}
}

type TypeRouter struct {
	galats.PesimisticLocking
}
