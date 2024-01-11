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

func (l *LiabilityRouter) Del(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	l.Lock()
	defer l.Done(w)
	var err error
	l.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	l.Tx, err = l.Db.Begin()
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Key not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(l.Db, "SELECT*FROM account WHERE tipe=? AND name=?", "Liability", name)
	if rows == nil {
		panic(fmt.Errorf("Liability not found"))
	}
	_, err = l.Db.Exec("DELETE FROM account WHERE name=? AND tipe=?", name, "Liability")
	if err != nil {
		panic(err)
	}
	if err = json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (l *LiabilityRouter) Get(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Key not found"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=? AND name=?", "Liability", name)
	if rows == nil {
		panic(fmt.Errorf("Liability not found"))
	}
	datas := arrayutils.Map(rows, row2Liability)
	one := datas[0]
	if err := json.NewEncoder(w).Encode(&one); err != nil {
		panic(err)
	}
}

func (l *LiabilityRouter) Add(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	l.Lock()
	defer l.Done(w)
	var err error
	l.Db, err = config.LoadDb()
	if err != nil {
		panic(err)
	}
	l.Tx, err = l.Db.Begin()
	if err != nil {
		panic(err)
	}
	var body models.Liability
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	if _, err = config.Str2Time(body.DueDate); err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(l.Db, "SELECT*FROM account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Liability %s is exists", body.Name))
	}
	if _, err = l.Db.Exec("INSERT INTO account(name,amount,interest,dueDate,tipe) VALUES(?,?,?,?,?)", body.Name, body.Amount, body.Interest, body.DueDate, "Liability"); err != nil {
		panic(err)
	}
	if err = json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (l *LiabilityRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=?", "Liability")
	if err := json.NewEncoder(w).Encode(&models.Liabilities{
		Datas: arrayutils.Map(rows, row2Liability),
	}); err != nil {
		panic(err)
	}
}

func row2Liability(v map[string]string, index int) models.Liability {
	amount, err := strconv.ParseFloat(v["amount"], 64)
	if err != nil {
		panic(err)
	}
	interest, err := strconv.ParseFloat(v["interest"], 64)
	if err != nil {
		panic(err)
	}
	return models.Liability{Name: v["name"], Amount: amount, Interest: interest, DueDate: v["dueDate"]}
}

func NewLiabilityRouter() LiabilityRouter {
	return LiabilityRouter{}
}

type LiabilityRouter struct {
	galats.PesimisticLocking
}
