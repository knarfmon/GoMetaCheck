package main

import "fmt"

func main() {
	fmt.Println("Hello Frankly Boy")

	foo()
	iter()
}

func foo() {
	fmt.Println("Hello again")
}

func iter() {
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}
}
