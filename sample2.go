package main

import (
	"fmt"
)

func main() {
	a := make([]int, 10)
	printSlice("a", a)
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s,
		len(x),
		cap(x),
		x)
}
