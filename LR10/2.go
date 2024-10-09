package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextHex string, key []byte) (string, error) {
	ciphertext, _ := hex.DecodeString(ciphertextHex)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("длина шифртекста слишком мала")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Использование: go run 2.go <mode> <string> <key>")
		fmt.Println("mode: encrypt или decrypt")
		return
	}

	mode := os.Args[1]
	data := os.Args[2]
	key := os.Args[3]

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Ключ должен быть длиной 16, 24 или 32 символа.")
		return
	}

	switch mode {
	case "encrypt":
		encrypted, err := encrypt(data, []byte(key))
		if err != nil {
			fmt.Printf("Ошибка при шифровании: %v\n", err)
			return
		}
		fmt.Printf("Зашифрованные данные: %s\n", encrypted)

	case "decrypt":
		decrypted, err := decrypt(data, []byte(key))
		if err != nil {
			fmt.Printf("Ошибка при расшифровании: %v\n", err)
			return
		}
		fmt.Printf("Расшифрованные данные: %s\n", decrypted)

	default:
		fmt.Println("Неизвестный режим. Используйте 'encrypt' или 'decrypt'.")
	}
}