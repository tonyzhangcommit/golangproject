package main

import "fmt"

type Person struct {
	age int
}

func main() {
	fmt.Println("hello world!")

	p := &Person{29}

	defer fmt.Println(p.age)
	
	defer func(p *Person) {
		fmt.Println(p.age)
	}(p)

	defer func() {
		fmt.Println(p.age)
	}()

	p = &Person{30}

}
