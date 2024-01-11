package services

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/galats"
	"main/models"
	"net/http"
	"strconv"

	arrayutils "github.com/AchmadRifai/array-utils"
	mapsutils "github.com/AchmadRifai/maps-utils"
	sqlutils "github.com/AchmadRifai/sql-utils"
	"github.com/gorilla/mux"
)

type AccountRouter struct {
	galats.PesimisticLocking
}

func (a *AccountRouter) Del(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	a.Lock()
	defer a.Done(w)
	var err error
	a.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	a.Tx, err = a.Db.Begin()
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Name param not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(a.Db, "SELECT*FROM account WHERE name=? AND tipe=?", name, "Account")
	if rows == nil {
		panic(fmt.Errorf("Data %s not found", name))
	}
	_, err = a.Db.Exec("DELETE FROM account WHERE name=? AND tipe=?", name, "Account")
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (a *AccountRouter) Add(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	a.Lock()
	defer a.Done(w)
	var err error
	a.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	a.Tx, err = a.Db.Begin()
	if err != nil {
		panic(err)
	}
	var body models.Account
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(a.Db, "SELECT*FROM account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Account %s is exists", body.Name))
	}
	_, err = a.Db.Exec("INSERT INTO account(name,amount,description,tipe) VALUES(?,?,?,?)", body.Name, body.Amount, body.Description, "Account")
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (a *AccountRouter) Get(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Name param not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE name=? AND tipe=?", name, "Account")
	if rows == nil {
		panic(fmt.Errorf("Data %s not found", name))
	}
	datas := arrayutils.Map(rows, func(v map[string]string, i int) models.Account {
		amount, err := strconv.ParseFloat(v["amount"], 64)
		if err != nil {
			panic(err)
		}
		return models.Account{Name: v["name"], Amount: amount, Description: v["description"]}
	})
	if err := json.NewEncoder(w).Encode(&datas[0]); err != nil {
		panic(err)
	}
}

func (a *AccountRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=?", "Account")
	if err := json.NewEncoder(w).Encode(&models.Accounts{
		Datas: arrayutils.Map(rows, func(v map[string]string, index int) models.Account {
			amount, err := strconv.ParseFloat(v["amount"], 64)
			if err != nil {
				panic(err)
			}
			return models.Account{Name: v["name"], Amount: amount, Description: v["description"]}
		}),
	}); err != nil {
		panic(err)
	}
}

func NewAccountRouter() AccountRouter {
	return AccountRouter{}
}

func initWriter(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}
