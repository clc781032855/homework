package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	for i := 1; i <= 10; i++ {
		fmt.Printf("生产者发送: %d\n", i)
		ch <- i // 发送数据到通道
	}
	fmt.Println("生产者完成")
}

func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // 完成时通知 WaitGroup

	// 使用 range 循环接收通道数据直到通道关闭
	for num := range ch {
		fmt.Printf("消费者接收: %d\n", num)
	}
	fmt.Println("消费者完成")
}

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go producer(ch, &wg)
	go consumer(ch, &wg)

	wg.Wait()
}
