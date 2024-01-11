package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitOutcome(r *mux.Router) {
	svc := services.NewOutcome()
	r.HandleFunc("/outcome", svc.Add).Methods(http.MethodPost)
	r.HandleFunc("/outcome/{name}", svc.Dell).Methods(http.MethodDelete)
	r.HandleFunc("/outcomes", svc.All).Methods(http.MethodGet)
}
