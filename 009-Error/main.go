package main

import "fmt"

func main() {

	n, err := fmt.Println("Hello, This is the error demonstration.")
	if err != nil {
		fmt.Println("Their is an error")
	}
	fmt.Println(n)
}
