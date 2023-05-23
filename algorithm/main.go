package main

import (
	"fmt"
)

type Math struct {
	x, y int
}

var m = map[string]Math{
	"foo": Math{2, 3},
}

func main() {
	v := []int{1,2,3}
	for i,value  := range v{
		fmt.Println((i))
		v = append(v, value)
	} 
	fmt.Println(v)

}
