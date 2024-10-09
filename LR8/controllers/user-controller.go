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

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
	Age  int                `bson:"age" json:"age"`
}

var userCollection *mongo.Collection

func InitUserController(collection *mongo.Database) {
	userCollection = collection.Collection("users")
}

func ValidateInput(user User) bool {
	if reflect.TypeOf(user.Name).String() == "string" ||
		user.Name != "" ||
		reflect.TypeOf(user.Age).String() == "int" ||
		(user.Age > 0 && user.Age < 150) {

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

	users := []User{}

	opt := options.Find().SetSkip(skip).SetLimit(int64(limit))
	cursor, err := userCollection.Find(ctx, filter, opt)

	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var usr User
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
	var usr User
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&usr)
	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var usr User
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
		handlers.HandleError(w, fmt.Errorf("incorrect value"), http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var usr User
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlers.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &usr)
	if err != nil || !ValidateInput(usr) {
		handlers.HandleError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": usr})
	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		handlers.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
