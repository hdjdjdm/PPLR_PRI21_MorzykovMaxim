package routes

import (
	"net/http"
	"server/controllers"

	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", controllers.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users", controllers.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", controllers.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", controllers.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", controllers.DeleteUser).Methods(http.MethodDelete)
}
