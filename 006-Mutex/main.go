package main

import (
	"sync"

	"fmt"
)






func main() {
	var wg sync.WaitGroup
	incrementor:= 0
	gs :=100
	wg.Add(gs)
	m:= sync.Mutex{}

	for i :=0; i<gs;i++{
		go func(){
			m.Lock()
			v:= incrementor
			//runtime.Gosched()
			v++
			incrementor = v
			fmt.Println(incrementor)
			m.Unlock()
			wg.Done()
		}()

	}

	wg.Wait()
	fmt.Println("end value:", incrementor)
}

