package routes

import (
	"net/http"
	"server/controllers"
	"server/middleware"

	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) {
	r.HandleFunc("/users/login", controllers.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/register", controllers.RegisterHandler).Methods(http.MethodPost)

	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware)

	userRoutes.HandleFunc("", controllers.GetUsers).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}", controllers.GetUser).Methods(http.MethodGet)

	adminRoutes := userRoutes.PathPrefix("").Subrouter()
	adminRoutes.Use(middleware.RoleMiddleware("admin"))
	adminRoutes.HandleFunc("", controllers.CreateUser).Methods(http.MethodPost)
	adminRoutes.HandleFunc("/{id}", controllers.UpdateUser).Methods(http.MethodPut)
	adminRoutes.HandleFunc("/{id}", controllers.DeleteUser).Methods(http.MethodDelete)

	userRoutes.HandleFunc("/admin-access", controllers.AdminAccessHandler).Methods(http.MethodPost)
}
