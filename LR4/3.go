package main

import "fmt"

func main() {
	var s string
	var people = map[string]int {
		"Макс": 20,
		"норМакс": 2,
		"ээээ": 100,
		"КРУТОЙ_ЧЕЛ": 666,
	}
	
	fmt.Println("Начальная карта:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
	
	fmt.Print("\nВведите имя для удаления: ")
	fmt.Scan(&s)
	for name, _ := range people {
		if name == s {
			delete(people, s)
			fmt.Printf("'%s' успешно удален\n", s)
			fmt.Println("\nОбновленная карта:")
			for name, age := range people {
				fmt.Printf("%s: %d лет\n", name, age)
			}
			return;
		}
	}
	fmt.Printf("Имя '%s' не найдено :/\n", s)
}