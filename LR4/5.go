package main

import "fmt"

func main() {
	var arr []int
	n, sum := 0, 0
	
	fmt.Print("Введите количество чисел: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	for i := 1; i <= n; i++ {
		var num int
		fmt.Printf("Введите число(%d/%d): ", i, n)
		_, err = fmt.Scan(&num)
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			return
		}
		arr = append(arr, num)
	}
	
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}

	fmt.Println("\nСумма всех чисел:", sum)
}