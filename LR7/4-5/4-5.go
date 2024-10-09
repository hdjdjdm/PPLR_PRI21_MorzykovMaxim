package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Начало: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		log.Printf("Завершено: %s %s за %v", r.Method, r.URL.Path, duration)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET запросы разрешены", http.StatusMethodNotAllowed)
		return
	}
	message := "Привет, мир!"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Полученные данные: Name = %s, Value = %s\n", data.Name, data.Value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные успешно получены"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)

	loggedMux := loggingMiddleware(mux)

	fmt.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
