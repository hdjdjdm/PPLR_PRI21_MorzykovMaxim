package mathutils

import "fmt"

func Factorial(n int) int {
	if n < 0 {
		fmt.Println("Ошибка: факториал отрицательного числа не определен.")
		return 0;
	}

	if n == 0 {
		return 1
	}

	return n * Factorial(n - 1)
}