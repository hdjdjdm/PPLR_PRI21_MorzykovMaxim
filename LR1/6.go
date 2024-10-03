package main

import "fmt"

func avg(a, b, c float64) float64 {
	sum := a + b + c
	return sum / 3
}

func main() {
	var a, b, c float64

	fmt.Print("Введите три числа: ")
	_, err := fmt.Scan(&a, &b, &c)
	if err != nil {
		fmt.Println("Ошибка ввода числа:", err)
		return
	}

	fmt.Printf("Среднее значение для %.2f, %.2f и %.2f = %.2f\n", a, b, c, avg(a, b, c))
}
