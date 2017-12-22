package main

import (
	"fmt"
	"sync"
	"runtime"
)

func main() {

	fmt.Println("begin cpu",runtime.NumCPU())
	fmt.Println("begin gr",runtime.NumGoroutine())

	var wg sync.WaitGroup
	wg.Add(2)

	func() { fmt.Println("Hello from Thing One")
	wg.Done()}()
	fmt.Println("begin cpu",runtime.NumCPU())
	fmt.Println("begin gr",runtime.NumGoroutine())
	func() { fmt.Println("Hello from Thing Two")
	wg.Done()}()

	fmt.Println("begin cpu",runtime.NumCPU())
	fmt.Println("begin gr",runtime.NumGoroutine())

	wg.Wait()

	fmt.Println("begin cpu",runtime.NumCPU())
	fmt.Println("begin gr",runtime.NumGoroutine())

	fmt.Println("About to exit")

}
