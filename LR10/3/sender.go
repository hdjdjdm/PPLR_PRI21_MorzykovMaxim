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

func signMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, privateKey, 0, message)
}

func encryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func decryptMessage(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func main() {
	privateKeySender, err := loadPrivateKey("sender_private_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	publicKeyReceiver, err := loadPublicKey("receiver_public_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("Сообщение от отправителя.")
	fmt.Println("Сообщение для отправки: 'Сообщение от отправителя'")

	signature, err := signMessage(privateKeySender, message)
	if err != nil {
		log.Fatal(err)
	}

	encryptedMessage, err := encryptMessage(publicKeyReceiver, message)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Отправляем сообщение...")
	_, err = conn.Write(append(signature, encryptedMessage...))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Сообщение отправлено.")

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	response, err := decryptMessage(privateKeySender, buffer[:n])
	if err != nil {
		log.Fatal("Ошибка при расшифровке ответа:", err)
	}

	fmt.Printf("Получено от получателя: %s\n", string(response))
}
