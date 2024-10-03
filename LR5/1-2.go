package main

import "fmt"

type Person struct {
	name string
	age int
}

func (p Person) getInfo() {
	fmt.Printf("Имя: %s\nВозраст: %d\n", p.name, p.age)
}

func (p *Person) birthday() {
	(*p).age++
}

func main() {
	p := Person{"КРУТОЙ_ЧЕЛ", 666}
	
	fmt.Println("До дня рождения: ")
	p.getInfo()
	
	p.birthday()
	fmt.Println("\nПосле дня рождения: ")
	p.getInfo()
}