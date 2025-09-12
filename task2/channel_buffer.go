package main

import (
	"fmt"
	"sync"
)

func producer_buffer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)
	for i := 0; i < 100; i++ {
		ch <- i
		fmt.Printf("producer正在发送整数:%d,目前channel长度:%d,channel容量:%d\n", i, len(ch), cap(ch))
	}
	fmt.Println("producer发送完成")
}

func consumer_buffer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("consumer正在接收整数:%d,目前channel长度:%d,channel容量:%d\n", num, len(ch), cap(ch))
	}
	fmt.Println("consumer接收完成")
}

func main() {
	ch := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)
	go producer_buffer(ch, &wg)
	go consumer_buffer(ch, &wg)
	wg.Wait()
}
