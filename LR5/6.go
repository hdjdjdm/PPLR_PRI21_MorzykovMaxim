package main

import "fmt"

type Stringer interface {
	showInfo()
}

type Book struct {
	title string
	author string
	year int
	price float64
}

func (b Book) showInfo() {
	fmt.Printf("Название: '%s'\n", b.title)
	fmt.Printf("Автор: %s\n", b.author)
	fmt.Printf("Год издания: %d\n", b.year)
	fmt.Printf("Цена: %.2f$\n", b.price)
}

func main() {
	book := Book{
		title: "Классная книжка",
		author: "Я",
		year: 2025,
		price: 24.99,
	}
	book.showInfo()
}