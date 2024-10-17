package main

import (
	"fmt"
	"net/http"
	"server/controllers"
	db "server/database"
	"server/routes"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Database

func main() {
	userCollection = db.DBConnect("mongodb://localhost:27017/", "LR8_DB")
	controllers.InitUserController(userCollection)
	controllers.InitLoginController(userCollection)

	r := mux.NewRouter()
	routes.SetUserRoutes(r)

	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", r)
}
