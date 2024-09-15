package main

import "fmt"

// CalcInt выполняет арифметические операции с двумя целыми числами.
func CalcInt(a, b int, op string) {
	switch op {
	case "+":
		fmt.Printf("%d + %d = %d\n", a, b, a+b)
	case "-":
		fmt.Printf("%d - %d = %d\n", a, b, a-b)
	case "*":
		fmt.Printf("%d * %d = %d\n", a, b, a*b)
	case "/":
		if b != 0 {
			fmt.Printf("%d / %d = %.2f\n", a, b, float64(a)/float64(b))
		} else {
			fmt.Println("Ошибка: Деление на 0 невозможно")
		}
	default:
		fmt.Println("Ошибка: Недоступная операция:", op)
	}
}

func main() {
	var a, b int
	var op string

	fmt.Println("Доступные операции: +, -, *, /")
	fmt.Print("Введите выражение (например '1 + 2'): ")
	_, err := fmt.Scan(&a, &op, &b)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	CalcInt(a, b, op)
}
