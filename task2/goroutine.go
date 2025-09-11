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

func main() {
	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	go oddnum(num)
	go evennum(num)
	time.Sleep(5000 * time.Millisecond)
}
