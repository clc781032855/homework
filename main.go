package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func singleNumber(nums []int) int {
	countmap := make(map[int]int)
	for _, num := range nums {
		countmap[num]++
	}
	for num, count := range countmap {
		if count == 1 {
			return num
		}
	}
	return 0
}

func main() {
	nums := []int{1, 1, 2, 3, 3}
	fmt.Println(singleNumber(nums))
}
