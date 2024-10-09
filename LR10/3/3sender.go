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

// Функция для подписи сообщения
func signMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, privateKey, 0, message)
}

// Функция для шифрования сообщения
func encryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func decryptMessage(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// Основная функция программы отправителя
func main() {
	// Загрузка закрытого ключа отправителя
	privateKeySender, err := loadPrivateKey("sender_private_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Загрузка открытого ключа получателя
	publicKeyReceiver, err := loadPublicKey("receiver_public_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Сообщение для отправки
	message := []byte("Сообщение от отправителя.")
	fmt.Println("Сообщение для отправки: 'Сообщение от отправителя'")

	// Подпись сообщения
	signature, err := signMessage(privateKeySender, message)
	if err != nil {
		log.Fatal(err)
	}

	// Шифрование сообщения
	encryptedMessage, err := encryptMessage(publicKeyReceiver, message)
	if err != nil {
		log.Fatal(err)
	}

	// Установка соединения с получателем
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Отправка сообщения, подписи и шифрованного сообщения
	fmt.Println("Отправляем сообщение...")
	_, err = conn.Write(append(signature, encryptedMessage...))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Сообщение отправлено.")

	// Ожидание зашифрованного ответа от получателя
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Расшифровка ответа
	response, err := decryptMessage(privateKeySender, buffer[:n])
	if err != nil {
		log.Fatal("Ошибка при расшифровке ответа:", err)
	}

	// Вывод расшифрованного ответа
	fmt.Printf("Получено от получателя: %s\n", string(response))
}
