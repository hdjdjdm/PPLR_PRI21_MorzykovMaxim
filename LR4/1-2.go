package main

import "fmt"

func avgAge(people map[string]int) float64 {
	sum := 0
	for _, age := range people {
		sum += age
	}

	return float64(sum / len(people))
}

func main() {
	var people = map[string]int {
		"Макс": 20,
		"норМакс": 2,
		"ээээ": 100,
	}
	
	fmt.Println("Начальная карта:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
	
	people["КРУТОЙ_ЧЕЛ"] = 666
	
	fmt.Println("\nОбновленная карта:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}

	fmt.Printf("\nСредний возраст: %.2f лет\n", avgAge(people))
}