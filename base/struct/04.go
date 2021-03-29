package main

import "fmt"

// 匿名结构体赋值
func main() {
	// 列表式的方式初始化结构体
	person := struct {
		Name  string
		Age   int
		Score float64
	}{
		"Forest",
		22,
		99.9,
	}

	// 键值对的方式初始化结构体
	student := struct {
		Name  string
		Age   int
		Score float64
	}{
		Name:  "Forest",
		Age:   22,
		Score: 99.9,
	}

	type Student struct {
		Name string
		Age int
		Score float64
	}
	

	fmt.Println("匿名结构体列表式赋值：", person)
	fmt.Println("匿名结构体键值对式赋值：", student)
}
