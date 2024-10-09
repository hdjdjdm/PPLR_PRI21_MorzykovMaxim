package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Message)
var mu sync.Mutex

type Message struct {
	ClientID  string `json:"client_id"`
	Text      string `json:"text"`
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении до WebSocket: %v", err)
		return
	}
	defer ws.Close()

	clientAddr := getClientAddress(ws)

	mu.Lock()
	clients[ws] = clientAddr
	mu.Unlock()

	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}

		broadcast <- Message{ClientID: clientAddr, Text: msg}
	}
}

func getClientAddress(ws *websocket.Conn) string {
	clientAddr, _, _ := net.SplitHostPort(ws.RemoteAddr().String())
	return clientAddr
}

func handleMessages() {
	for {
		msg := <-broadcast

		mu.Lock()
		for client, addr := range clients {
			if addr != "" {
				err := client.WriteJSON(fmt.Sprintf("[%s] отправил сообщение: %s", msg.ClientID, msg.Text))
				if err != nil {
					log.Printf("Ошибка при отправке сообщения: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}
