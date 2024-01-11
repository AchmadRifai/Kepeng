package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitLiability(r *mux.Router) {
	svc := services.NewLiabilityRouter()
	r.HandleFunc("/liability", svc.Add).Methods(http.MethodPost)
	r.HandleFunc("/liabilities", svc.All).Methods(http.MethodGet)
	r.HandleFunc("/liability/{name}", svc.Get).Methods(http.MethodGet)
	r.HandleFunc("/liability/{name}", svc.Del).Methods(http.MethodDelete)
}
