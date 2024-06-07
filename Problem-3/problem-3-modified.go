package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	cnp := make(chan func(), 10)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range cnp {
				f()
			}
		}()
	}

	cnp <- func() {
		fmt.Println("HERE1")
	}

	close(cnp)

	wg.Wait()
	fmt.Println("Hello")
}
