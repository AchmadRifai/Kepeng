package routers

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func InitType(r *mux.Router) {
	svc := services.NewTypeRouter()
	r.HandleFunc("/types", svc.All).Methods(http.MethodGet)
	r.HandleFunc("/type", svc.Add).Methods(http.MethodPost)
	r.HandleFunc("/type/{name}", svc.Del).Methods(http.MethodDelete)
}
