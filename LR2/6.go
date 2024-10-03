package main

import (
	"fmt"
)

func avg(a, b int) float64 {
	sum := float64(a + b)
	return sum / 2
}

func main() {
	var a, b int

	fmt.Print("Введите два целых числа: ")
	_, err := fmt.Scan(&a, &b)
	if err != nil {
		fmt.Println("Ошибка ввода числа:", err)
		return
	}

	fmt.Printf("Среднее значение для %d и %d = %.2f\n", a, b, avg(a, b))
}
