package main

import (
	"fmt"
)

func main() {
	type Person struct {
		Name   string
		Age    int
		Gender int
		Sign   string
	}

	type Developer struct {
		Person
		Language string
		skill    string
	}

	// 列表式初始化结构体
	developer := Developer{
		Person{
			"Forset",
			20,
			0,
			"君子不器",
		},
		"javascript、nodejs、go",
		"javascript、nodejs、go",
	}

	// 键值对的方式初始化结构体
	developer1 := Developer{
		Person: Person{
			Name: "Forest",
			Age:  22,
		},
		Language: "javascript css nodejs go",
		skill:    "vue react mysql mongoDB redis docker nginx 微信小程序 微信公众号",
	}

	fmt.Println("列表式初始化结构体:", developer)     // 列表式初始化结构体: {{Forset 20 0 君子不器} javascript、nodejs、go javascript、nodejs、go}
	fmt.Println("键值对的方式初始化结构体:", developer1) // 键值对的方式初始化结构体: {{Forest 22 0 } javascript css nodejs go vue react mysql mongoDB redis docker nginx 微信小程序 微信公众号}
}
