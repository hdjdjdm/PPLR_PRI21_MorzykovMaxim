package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"strings"
)

func computeHash(data, algorithm string) (string, error) {
	var h hash.Hash
	switch strings.ToLower(algorithm) {
	case "md5":
		h = md5.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("неподдерживаемый алгоритм хэширования: %s", algorithm)
	}

	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func verifyIntegrity(data, expectedHash, algorithm string) (bool, error) {
	hash, err := computeHash(data, algorithm)
	if err != nil {
		return false, err
	}
	return hash == expectedHash, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: go run 1.go <строка> <алгоритм> <ожидаемый_хэш (необязательно)>")
		return
	}

	data := os.Args[1]
	algorithm := os.Args[2]

	if len(os.Args) == 4 {
		expectedHash := os.Args[3]
		match, err := verifyIntegrity(data, expectedHash, algorithm)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		if match {
			fmt.Println("Хэш совпадает с ожидаемым.")
		} else {
			fmt.Println("Хэш не совпадает с ожидаемым.")
		}
	} else {
		hash, err := computeHash(data, algorithm)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		fmt.Printf("Вычисленный %s хэш: %s\n", strings.ToUpper(algorithm), hash)
	}
}
