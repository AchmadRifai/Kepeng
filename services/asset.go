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

func (a *AssetRouter) Del(w http.ResponseWriter, r *http.Request) {
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
		panic(fmt.Errorf("Key name not exists"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(a.Db, "SELECT*FROM account WHERE name=?", name)
	if rows == nil {
		panic(fmt.Errorf("asset %s not found", name))
	}
	_, err = a.Db.Exec("DELETE FROM account WHERE name=?", name)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (a *AssetRouter) Add(w http.ResponseWriter, r *http.Request) {
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
	var body models.Asset
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(a.Db, "SELECT*FROM account WHERE name=?", body.Name)
	if rows != nil {
		panic(fmt.Errorf("Asset %s is exists", body.Name))
	}
	_, err = a.Db.Exec("INSERT INTO account(name,amount,increase,decrease,tipe) VALUES(?,?,?,?,?)", body.Name, body.Amount, body.Increase, body.Decrease, "Asset")
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(&models.MessageRes{Message: "Success"}); err != nil {
		panic(err)
	}
}

func (a *AssetRouter) Get(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	if !arrayutils.Contains(mapsutils.KeysOfMap(params), "name") {
		panic(fmt.Errorf("Key name not exists"))
	}
	name := params["name"]
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE name=? AND tipe=?", name, "Asset")
	if rows == nil {
		panic(fmt.Errorf("asset %s not found", name))
	}
	datas := arrayutils.Map(rows, rowToAsset)
	data := datas[0]
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		panic(err)
	}
}

func (a *AssetRouter) All(w http.ResponseWriter, r *http.Request) {
	defer galats.NormalErrors()
	initWriter(w)
	db, err := config.LoadDb()
	defer galats.Db500(w, db)
	if err != nil {
		panic(err)
	}
	rows, _ := sqlutils.DbSelect(db, "SELECT*FROM account WHERE tipe=?", "Asset")
	if err := json.NewEncoder(w).Encode(&models.Assets{
		Datas: arrayutils.Map(rows, rowToAsset),
	}); err != nil {
		panic(err)
	}
}

func rowToAsset(v map[string]string, index int) models.Asset {
	amount, err := strconv.ParseFloat(v["amount"], 64)
	if err != nil {
		panic(err)
	}
	increase, err := strconv.ParseFloat(v["increase"], 64)
	if err != nil {
		panic(err)
	}
	decrease, err := strconv.ParseFloat(v["decrease"], 64)
	if err != nil {
		panic(err)
	}
	return models.Asset{Name: v["name"], Amount: amount, Increase: increase, Decrease: decrease}
}

func NewAssetRouter() AssetRouter {
	return AssetRouter{}
}

type AssetRouter struct {
	galats.PesimisticLocking
}
