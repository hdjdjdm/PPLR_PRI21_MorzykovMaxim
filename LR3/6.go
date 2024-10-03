package main

import (
	"fmt"
)

func main() {
	var strings []string
	n, maxlen, maxindex := 0, 0, 0
	
	fmt.Print("Введите количество строк: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	for i := 0; i < n; i++ {
		var str string
		fmt.Printf("Введите строку(%d/%d): ", i, n)
		fmt.Scan(&str)
		strings = append(strings, str)
	}
	
	for i := 0; i < len(strings); i++ {
		if (len(strings[i]) > maxlen) {
			maxlen = len(strings[i])
			maxindex = i
		}
	}

	fmt.Println("Самое длинное слово:", strings[maxindex])
	fmt.Println("Длина:", maxlen)
}