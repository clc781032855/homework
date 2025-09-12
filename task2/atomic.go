package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	const (
		numGoroutines = 10
		numIncrements = 1000
	)

	var counter int64

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numIncrements; j++ {

				atomic.AddInt64(&counter, 1)
			}
		}(i)
	}

	wg.Wait()

	finalValue := atomic.LoadInt64(&counter)

	fmt.Printf("最终计数器值: %d (预期: %d)\n", finalValue, numGoroutines*numIncrements)

	if finalValue != int64(numGoroutines*numIncrements) {
		fmt.Printf("错误: 实际值 %d != 预期值 %d\n", finalValue, numGoroutines*numIncrements)
	} else {
		fmt.Println("验证成功: 计数器值正确")
	}
}
