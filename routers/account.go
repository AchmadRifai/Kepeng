package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAccount(r *mux.Router) {
	service := services.NewAccountRouter()
	r.HandleFunc("/account", service.Add).Methods(http.MethodPost)
	r.HandleFunc("/account/{name}", service.Get).Methods(http.MethodGet)
	r.HandleFunc("/accounts", service.All).Methods(http.MethodGet)
	r.HandleFunc("/account/{name}", service.Del).Methods(http.MethodDelete)
}
