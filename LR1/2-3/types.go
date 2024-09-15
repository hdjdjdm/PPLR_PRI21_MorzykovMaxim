package main

import "fmt"

func main() {
	i := 1
	f := 1.1
	s := "Hewwoo:з"
	b := true

	fmt.Println("Типы переменных")
	fmt.Printf("i = %d,\tтип: %T\n", i, i)
	fmt.Printf("f = %.2f,\tтип: %T\n", f, f)
	fmt.Printf("s = %s,\tтип: %T\n", s, s)
	fmt.Printf("b = %t,\tтип: %T\n", b, b)
}
