package main

import "fmt"

func length(s string) int {
	len := 0
	for range s {
		len++
	}
	return len
	// return len(s)
}

func main() {
	var s string

	fmt.Print("Введите строку: ")
	_, err := fmt.Scan(&s)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	fmt.Printf("Длина строки: %d\n", length(s))
}