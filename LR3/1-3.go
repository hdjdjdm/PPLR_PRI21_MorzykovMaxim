package main

import (
	"fmt"

	"lr3/mathutils"
	"lr3/stringutils"
)

func main() {
	// 2
	var n int

	fmt.Print("Введите число: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	fmt.Printf("Факториал числа: %d\n", mathutils.Factorial(n))

	// 3
	var s string

	fmt.Print("Введите строку: ")
	_, err = fmt.Scan(&s)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	fmt.Printf("Перевернутая строка: %s\n", stringutils.Reverse(s))

}