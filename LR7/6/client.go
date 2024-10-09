package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8080/ws"

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()
	fmt.Println("Вы подключились к чату. Введите 'exit' для закрытия")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	go func() {
		for {
			var msg string
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("Ошибка при чтении сообщения: %v", err)
				return
			}
			fmt.Println(msg)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		err := conn.WriteJSON(text)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			return
		}
	}

	<-done
	log.Println("Получен сигнал завершения, закрытие подключения...")
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(time.Second)
}
