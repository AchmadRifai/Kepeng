package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitIncome(r *mux.Router) {
	svc := services.NewIncomeRouter()
	r.HandleFunc("/income", svc.Add).Methods(http.MethodPost)
	r.HandleFunc("/income/{name}", svc.Del).Methods(http.MethodDelete)
	r.HandleFunc("/incomes", svc.All).Methods(http.MethodGet)
}
