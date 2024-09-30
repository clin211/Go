package main

import "fmt"

func main() {
	s1 := "this is a string"

	fmt.Println(s1)

	i := 1
	fmt.Println(i != 0 || i != 1)

	arr := []int{1, 2, 3}
	for _, i := range arr {
		go func() {
			fmt.Println(i)
		}()
	}

}
