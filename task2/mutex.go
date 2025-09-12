package main

import (
	"fmt"
	"sync"
)

type Sharednum struct {
	count int
	mu    sync.Mutex
}

func (sn *Sharednum) AutoIncrement() {
	sn.mu.Lock()
	defer sn.mu.Unlock()
	sn.count++
}

func (sn *Sharednum) GetValue() int {
	sn.mu.Lock()
	defer sn.mu.Unlock()
	return sn.count
}

func main() {
	count := Sharednum{}

	const (
		numgoroutine = 10
		numincrement = 1000
	)

	var wg sync.WaitGroup

	wg.Add(numgoroutine)

	for i := 0; i < numgoroutine; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numincrement; j++ {
				count.AutoIncrement()
			}
			fmt.Printf("goroutine %d: %d\n", id, count.GetValue())
		}(i)
	}

	wg.Wait()

}
