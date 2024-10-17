package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"server/handlers"
	"server/models" // Импортируйте пакет models

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection

func InitUserController(collection *mongo.Database) {
	userCollection = collection.Collection("users")
}

func ValidateInput(user models.User) bool { // Обновите параметр функции
	if reflect.TypeOf(user.Name).String() == "string" ||
		user.Name != "" ||
		reflect.TypeOf(user.Age).String() == "int" || user.Age > 0 {

		return true
	}
	return false
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	nameFilter := r.URL.Query().Get("name")
	ageFilterStr := r.URL.Query().Get("age")

	page := 1
	limit := 10
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}
	}

	filter := bson.M{}
	if nameFilter != "" {
		filter["name"] = bson.M{"$regex": nameFilter, "$options": "i"}
	}
	if ageFilterStr != "" {
		age, err := strconv.Atoi(ageFilterStr)
		if err == nil {
			filter["age"] = age
		}
	}
	var skip int64 = int64((page - 1) * limit)

	users := []models.User{} // Измените тип на models.User

	opt := options.Find().SetSkip(skip).SetLimit(int64(limit))
	cursor, err := userCollection.Find(ctx, filter, opt)

	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var usr models.User // Используйте models.User
		if err := cursor.Decode(&usr); err != nil {
			handlers.HandleError(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, usr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var usr models.User // Используйте models.User
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handlers.HandleError(w, fmt.Errorf("Неверный формат идентификатора"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handlers.HandleError(w, fmt.Errorf("Пользователь не найден"), http.StatusNotFound)
			return
		}
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var usr models.User // Используйте models.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &usr)
	if err != nil {
		handlers.HandleError(w, err, http.StatusBadRequest)
		return
	}

	if !ValidateInput(usr) {
		handlers.HandleError(w, fmt.Errorf("Неверные данные пользователя"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usr.Id = primitive.NewObjectID()
	_, err = userCollection.InsertOne(ctx, usr)
	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "Пользователь успешно создан.",
		"user":    usr,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var usr models.User
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handlers.HandleError(w, fmt.Errorf("Неверный формат идентификатора"), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &usr)
	if err != nil || !ValidateInput(usr) {
		handlers.HandleError(w, fmt.Errorf("Неверные данные пользователя"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name": usr.Name,
			"age":  usr.Age,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Пользователь успешно обновлен.",
		"user":    usr,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handlers.HandleError(w, fmt.Errorf("Неверный формат идентификатора"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		handlers.HandleError(w, fmt.Errorf("Ошибка при удалении пользователя"), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		handlers.HandleError(w, fmt.Errorf("Пользователь не найден"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно удален"})
}
