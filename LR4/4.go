package main

import (
	"fmt"
	"strings"
)

func main() {
	var s string
	
	fmt.Print("Введите строку: ")
	fmt.Scan(&s)
	
	fmt.Printf("Новая строка: %s\n", strings.ToUpper(s))
}