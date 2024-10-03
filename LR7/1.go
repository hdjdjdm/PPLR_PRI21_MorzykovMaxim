package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup
var shutdown = false

func main() {
	// Создание TCP слушателя на порту 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка при создании слушателя:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Сервер запущен на порту 8080")

	// Канал для обработки сигналов завершения
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для обработки сигналов завершения
	go func() {
		<-stopChan
		fmt.Println("\nПолучен сигнал остановки. Ожидание завершения соединений...")
		shutdown = true
		ln.Close() // Закрываем слушатель, чтобы выйти из цикла Accept
		wg.Wait()  // Ожидаем завершения всех соединений
		fmt.Println("Сервер остановлен")
		os.Exit(0)
	}()

	// Основной цикл для принятия соединений
	for {
		if shutdown {
			break
		}

		conn, err := ln.Accept()
		if err != nil {
			if shutdown {
				break
			}
			fmt.Println("Ошибка при приеме соединения:", err)
			continue
		}

		wg.Add(1) // Увеличиваем счетчик ожидания
		go handleConnection(conn) // Обрабатываем соединение в отдельной горутине
	}
}

// Функция для обработки соединения
func handleConnection(conn net.Conn) {
	defer wg.Done() // Уменьшаем счетчик ожидания
	defer conn.Close()

	fmt.Println("Новое соединение", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Клиент закрыл соединение:", conn.RemoteAddr().String())
				return
			}
			fmt.Println("Ошибка при чтении:", err)
			break
		}

		message := string(buffer[:n])
		fmt.Printf("Получено сообщение от клиента: %s\n", message)

		response := "Сообщение получено\n"
		conn.Write([]byte(response))
	}
}
