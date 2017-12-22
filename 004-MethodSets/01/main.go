package main

import "fmt"

type person struct {
	fName string
}
type human interface {
	speak()
}

func (p *person) speak() {
	fmt.Println("Hello, I am a person")
}

func saySomething(h human) {
	h.speak()
}

func main() {

	p1 := person{
		fName: "Frank",
	}

	saySomething(&p1)

}
