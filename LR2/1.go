package main

import "fmt"

func main() {
	var a int

	fmt.Print("Введите число: ")
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	if (a % 2 == 0) {
		fmt.Print("Число четное")
	} else {
		fmt.Print("Число нечентное")
	}
}