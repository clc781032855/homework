package main

import (
	"fmt"
	"strconv"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

// 136.只出现一次的数字
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

// 回文数
func isPalindrome(x int) bool {
	str := strconv.Itoa(x)
	for i := 0; i < len(str); i++ {
		if str[i] != str[len(str)-1-i] {
			return false
		}
	}
	return true
}

func main() {
	nums := []int{1, 1, 2, 3, 3}
	fmt.Println(singleNumber(nums))
	fmt.Println(isPalindrome(121))
}
