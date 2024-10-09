package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

// Функция для загрузки закрытого ключа
func loadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// Функция для загрузки открытого ключа
func loadPublicKey(filename string) (*rsa.PublicKey, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// Функция для проверки подписи
func verifySignature(publicKey *rsa.PublicKey, message, signature []byte) error {
	return rsa.VerifyPKCS1v15(publicKey, 0, message, signature)
}

// Функция для расшифровки сообщения
func decryptMessage(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// Основная функция программы получателя
func main() {
	// Загрузка закрытого ключа получателя
	privateKeyReceiver, err := loadPrivateKey("receiver_private_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Загрузка открытого ключа отправителя
	publicKeySender, err := loadPublicKey("sender_public_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Создание TCP-сервера
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Ожидание соединения...")

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Соединение установлено.")

	// Получение сообщения и подписи
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Разделение подписи и зашифрованного сообщения
	signature := buffer[:256]
	encryptedMessage := buffer[256:n]

	// Расшифровка сообщения
	message, err := decryptMessage(privateKeyReceiver, encryptedMessage)
	if err != nil {
		log.Fatal("Ошибка при расшифровке сообщения:", err)
	}

	// Проверка подписи
	err = verifySignature(publicKeySender, message, signature)
	if err != nil {
		log.Fatal("Подпись неверна.")
	}

	fmt.Printf("Получено сообщение: %s\n", string(message))
	fmt.Println("Подпись подтверждена.")

	// Подготовка ответа
	responseMessage := []byte("Ответ от получателя.")
	fmt.Println("Сообщение для отправки: 'Ответ от получателя'")
	// Шифрование ответа
	encryptedResponse, err := rsa.EncryptPKCS1v15(rand.Reader, publicKeySender, responseMessage)
	if err != nil {
		log.Fatal("Ошибка при шифровании ответа:", err)
	}

	// Отправка зашифрованного ответа
	_, err = conn.Write(encryptedResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ответ отправлен.")
}
