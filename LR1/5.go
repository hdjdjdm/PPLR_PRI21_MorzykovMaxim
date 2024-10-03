package main

import "fmt"

func CalcFloat(a, b float64, op string) {
	switch op {
	case "+":
		fmt.Printf("%.2f + %.2f = %.2f\n", a, b, a+b)
	case "-":
		fmt.Printf("%.2f - %.2f = %.2f\n", a, b, a-b)
	default:
		fmt.Println("Недоступная операция: ", op)
	}
}

func main() {
	var a, b float64
	op := ""

	fmt.Println("Доступные операции: +, -")
	fmt.Print("Введите выражение (например '1.8 + 2.5'): ")
	_, err := fmt.Scan(&a, &op, &b)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	CalcFloat(a, b, op)
}
