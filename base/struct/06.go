package main

import (
	"fmt"
)

func main() {
	type Person struct {
		Name string
		Age  int
	}
	// 整体列表式赋值
	var p Person
	p = Person{
		"Forest",
		22,
	}

	fmt.Println("整体列表式赋值:", p) // 整体列表式赋值: {Forest 22}

	// 结构体键值对赋值
	var stu Person = Person{
		Name: "lin",
		Age:  22,
	}

	fmt.Println("结构体键值对赋值:", stu) // 结构体键值对赋值: {lin 22}

	// 结构体字段单独赋值
	var developer Person
	developer.Name = "Forest"
	developer.Age = 22
	fmt.Println("结构体字段单独赋值:", developer) // 结构体字段单独赋值: {Forest 22}
}
