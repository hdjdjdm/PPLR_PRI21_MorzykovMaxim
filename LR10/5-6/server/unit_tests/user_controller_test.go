package unit_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"server/controllers"
	"server/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var testDB *mongo.Database
var idStr string = "60b6c8f1f1e2b1c3d4e5f6a7"

func TestMain(m *testing.M) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}

	testDB = client.Database("test_db")
	dropCollections(testDB)

	controllers.InitUserController(testDB)

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Fatalf("Ошибка создания ObjectID из строки: %v", err)
	}
	_, err = testDB.Collection("users").InsertOne(context.TODO(), models.User{
		Id:   objectID,
		Name: "SUPER-TEST",
		Age:  55,
	})
	if err != nil {
		log.Fatalf("Ошибка вставки тестового пользователя: %v", err)
	}
	code := m.Run()

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal("Ошибка отключения от MongoDB:", err)
	}

	os.Exit(code)
}

func dropCollections(db *mongo.Database) {
	collections, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal("Ошибка получения коллекций:", err)
	}
	for _, collection := range collections {
		if err := db.Collection(collection).Drop(context.TODO()); err != nil {
			log.Fatalf("Ошибка удаления коллекции %s: %v", collection, err)
		}
	}
}

func TestCreateUser(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	for _, val := range testUsers {
		body, _ := json.Marshal(val)

		req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		if err != nil {
			t.Fatal("Ошибка создания запроса:", err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Неверный код статуса: получен %v, ожидается %v",
				status, http.StatusCreated)
		}
	}
}

func TestGetUsers(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")

	req, err := http.NewRequest("GET", "/users?page=1&limit=10", nil)
	if err != nil {
		t.Fatal("Ошибка создания запроса:", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код статуса: получен %v, ожидается %v",
			status, http.StatusOK)
	}

	var users []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	fmt.Println("Количество пользователей:", len(users))
	if len(users) != len(testUsers) {
		t.Errorf("Ожидалось %d пользователей, получено %d", len(testUsers), len(users))
	}
}

func TestUpdateUser(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")

	user := models.User{Name: "Jane Doe", Age: 25}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("PUT", "/users/60b6c8f1f1e2b1c3d4e5f6a7", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Ошибка создания запроса:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код статуса: получен %v, ожидается %v",
			status, http.StatusOK)
	}
}

func TestDeleteUser(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/users/60b6c8f1f1e2b1c3d4e5f6a7", nil)
	if err != nil {
		t.Fatal("Ошибка создания запроса:", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код статуса: получен %v, ожидается %v",
			status, http.StatusOK)
	}

	expected := map[string]string{"message": "Пользователь успешно удален"}
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["message"] != expected["message"] {
		t.Errorf("Неверное тело ответа: получено %v, ожидается %v",
			response, expected)
	}
}
