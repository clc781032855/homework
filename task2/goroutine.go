package main

import (
	"fmt"
	"time"
)

func oddnum(num []int) {
	for i := range num {
		if num[i]%2 != 0 {
			fmt.Printf("%d ", num[i])
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func evennum(num []int) {
	for i := range num {
		if num[i]%2 == 0 && num[i] > 1 {
			fmt.Printf("%d ", num[i])
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func task1() {
	start := time.Now()
	fmt.Printf("executing task1\n")
	time.Sleep(2000 * time.Millisecond)
	end := time.Now()
	fmt.Printf("task1 cost %v\n", end.Sub(start))
}
func task2() {
	start := time.Now()
	fmt.Printf("executing task2\n")
	time.Sleep(1000 * time.Millisecond)
	end := time.Now()
	fmt.Printf("task2 cost %v\n", end.Sub(start))
}
func task3() {
	start := time.Now()
	fmt.Printf("executing task3\n")
	time.Sleep(1500 * time.Millisecond)
	end := time.Now()
	fmt.Printf("task3 cost %v\n", end.Sub(start))
}

func main() {
	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	go oddnum(num)
	go evennum(num)
	go task1()
	go task2()
	go task3()
	time.Sleep(10000 * time.Millisecond)
}
