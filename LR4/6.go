package main

import "fmt"

func main() {
	var arr []int
	n := 0
	
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

	fmt.Println("\nНачальные числа:", arr)

	fmt.Print("\nЧисла в обратном порядке: ")
	for i := n-1; i >= 0; i-- {
		fmt.Printf("%d ", arr[i])
	}
}