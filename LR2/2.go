package main

import "fmt"

func pnz(a int) string {
	if (a > 0) {
		return "Positive"
	} else if (a < 0) {
		return "Negative"
	} else {
		return "Zero"
	}
}

func main() {
	var a int

	fmt.Print("Введите число: ")
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	fmt.Println(pnz(a))
}