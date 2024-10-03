package main

import "fmt"

func main() {
	var i int = 1
	var f float64 = 1.1
	var s string = "Hewwoo:з"
	var b bool = true
	
	fmt.Println("Типы переменных")
	fmt.Printf("i = %d,\tтип: %T\n", i, i)
	fmt.Printf("f = %.2f,\tтип: %T\n", f, f)
	fmt.Printf("s = %s,\tтип: %T\n", s, s)
	fmt.Printf("b = %t,\tтип: %T\n", b, b)
}
