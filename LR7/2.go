package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите сообщение: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении сообщения:", err)
		return
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Printf("Ответ от сервера: %s\n", string(buffer[:n]))
}
