package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"server/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("haha_my_secret_key")
var loginCollection *mongo.Collection
var secretKey = "tastycode__"

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InitLoginController(db *mongo.Database) {
	loginCollection = db.Collection("logins")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	var result struct {
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err = loginCollection.FindOne(r.Context(), bson.M{"login": creds.Login}).Decode(&result)
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.Claims{
		Login: creds.Login,
		Role:  result.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	var result struct {
		Login string `json:"login"`
	}
	err = loginCollection.FindOne(r.Context(), bson.M{"login": creds.Login}).Decode(&result)
	if err == nil {
		http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	newUser := bson.M{
		"login":    creds.Login,
		"password": string(hashedPassword),
		"role":     "user",
	}

	_, err = loginCollection.InsertOne(r.Context(), newUser)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация прошла успешно"})
}

func AdminAccessHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Secret string `json:"secret"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "неверный JSON", http.StatusBadRequest)
		return
	}

	if input.Secret != secretKey {
		http.Error(w, "неправильный код ~_~", http.StatusUnauthorized)
		return
	}

	currentUser, ok := r.Context().Value("claims").(*models.Claims)
	if !ok || currentUser == nil {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}

	if currentUser.Role == "admin" {
		http.Error(w, "вы уже обладаете правами администратора", http.StatusForbidden)
		return
	}

	currentUser.Role = "admin"
	_, err = loginCollection.UpdateOne(r.Context(), bson.M{"login": currentUser.Login}, bson.M{"$set": bson.M{"role": "admin"}})
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	newClaims := &models.Claims{
		Login: currentUser.Login,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "права администратора получены", "token": tokenString})
}
