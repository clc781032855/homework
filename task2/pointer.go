package main

import "fmt"

func changenumadd10(a *int) int {
	b := *a + 10
	return b
}

func doublenum(slice *[]int) {
	for i := range *slice {
		(*slice)[i] *= 2
	}
}

func main() {
	a := 10
	b := changenumadd10(&a)
	fmt.Println(a)
	fmt.Println(b)

	c := []int{1, 2, 3, 4, 5, 6}
	doublenum(&c)
	fmt.Println(c)
}
