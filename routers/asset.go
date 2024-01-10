package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAssetRouter(r *mux.Router) {
	svc := services.NewAssetRouter()
	r.HandleFunc("/asset", svc.Add).Methods(http.MethodPost)
	r.HandleFunc("/assets", svc.All).Methods(http.MethodGet)
	r.HandleFunc("/asset/{name}", svc.Del).Methods(http.MethodDelete)
	r.HandleFunc("/asset/{name}", svc.Get).Methods(http.MethodGet)
}
