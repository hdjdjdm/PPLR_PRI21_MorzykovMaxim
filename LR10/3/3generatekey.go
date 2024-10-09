package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func generateKeyPair() (*rsa.PrivateKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, err
    }

    return privateKey, nil
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
    out, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer out.Close()

    pem.Encode(out, &pem.Block{
        Type:  "PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })

    return nil
}

func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
    pubKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
    out, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer out.Close()

    pem.Encode(out, &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubKeyBytes,
    })

    return nil
}

func main() {
    senderPrivateKey, err := generateKeyPair()
    if err != nil {
        log.Fatal(err)
    }

    err = savePrivateKeyToFile(senderPrivateKey, "sender_private_key.pem")
    if err != nil {
        log.Fatal(err)
    }

    err = savePublicKeyToFile(&senderPrivateKey.PublicKey, "sender_public_key.pem")
    if err != nil {
        log.Fatal(err)
    }

    receiverPrivateKey, err := generateKeyPair()
    if err != nil {
        log.Fatal(err)
    }

    err = savePrivateKeyToFile(receiverPrivateKey, "receiver_private_key.pem")
    if err != nil {
        log.Fatal(err)
    }

    err = savePublicKeyToFile(&receiverPrivateKey.PublicKey, "receiver_public_key.pem")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Ключи успешно созданы и сохранены в файлы.")
}
