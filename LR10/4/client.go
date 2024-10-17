package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		fmt.Println("Ошибка при загрузке сертификата клиента:", err)
		os.Exit(1)
	}

	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Ошибка при чтении сертификата сервера (CA):", err)
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", "localhost:8080", config)
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
