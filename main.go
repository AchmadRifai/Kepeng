package main

import (
	"log"
	"main/config"
	"main/routers"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	config.InitDb()
	a.Router = mux.NewRouter()
	a.Router.Use(mux.CORSMethodMiddleware(a.Router))
	routers.InitAccount(a.Router)
	routers.InitAssetRouter(a.Router)
	routers.InitIncome(a.Router)
	routers.InitLiability(a.Router)
	routers.InitType(a.Router)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {
	a := App{}
	a.Initialize()
	a.Run(":8000")
}
