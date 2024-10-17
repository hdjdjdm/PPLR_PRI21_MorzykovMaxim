package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func main() {
	// Загрузка сертификатов сервера
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Ошибка при загрузке сертификата сервера:", err)
		os.Exit(1)
	}

	// Настройка пула доверенных сертификатов для клиентов (используем CA)
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Ошибка при чтении CA сертификата:", err)
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройки TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Взаимная аутентификация
	}

	ln, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Println("Ошибка при создании TLS слушателя:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("TLS сервер запущен на порту 8080")

	ctx, cancel := context.WithCancel(context.Background())

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println("\nПолучен сигнал остановки. Ожидание завершения соединений...")
		cancel()
		ln.Close()
		wg.Wait()
		fmt.Println("Сервер остановлен")
		os.Exit(0)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Завершаем приём новых соединений")
				return
			default:
				fmt.Println("Ошибка при приеме соединения:", err)
				continue
			}
		}

		wg.Add(1)
		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer wg.Done()

	done := make(chan struct{})
	go func() {
		defer conn.Close()
		defer close(done)

		fmt.Println("Новое TLS соединение", conn.RemoteAddr().String())
		buffer := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершаем соединение:", conn.RemoteAddr().String())
				return
			default:
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))

				n, err := conn.Read(buffer)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					if err.Error() == "EOF" {
						fmt.Println("Клиент закрыл соединение:", conn.RemoteAddr().String())
						return
					}
					fmt.Println("Ошибка при чтении:", err)
					return
				}

				message := string(buffer[:n])
				fmt.Printf("Получено сообщение от клиента: %s", message)

				response := "Сообщение получено\n"
				conn.Write([]byte(response))
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Принудительное завершение соединения", conn.RemoteAddr().String())
		conn.Close()
	case <-done:
	}
}
