package main

import (
	"fmt"
)

func main() {

	type Person struct {
		Name   string
		Age    int
		Gender int
	}

	type Developer struct {
		Person
		Language string
		skill    string
	}

	// 结构体匿名字段整体键值赋值
	developer := Developer{
		Person: Person{
			Name: "Forest",
			Age:  22,
		},
		Language: "javascript nodejs css go",
	}

	// 整体列表式赋值
	developer1 := Developer{
		Person{
			"Forest",
			22,
			0,
		},
		"javascript nodejs css go",
		"vue react docker",
	}

	// 单个字段赋值
	developer.skill = "append(developer.skill, \"vue\")"
	developer.skill = "append(developer.skill, \"react\")"
	developer.skill = "append(developer.skill, \"node\")"
	developer.skill = "append(developer.skill, \"docker\")"
	developer.skill = "append(developer.skill, \"MongoDB\")"

	fmt.Println("结构体匿名字段整体键值赋值和单独赋值:", developer) // 结构体匿名字段整体键值赋值和单独赋值: {{Forest 22 0} javascript nodejs css go append(developer.skill, "MongoDB")}
	fmt.Println("整体列表式赋值:", developer1)           // 整体列表式赋值: {{Forest 22 0} javascript nodejs css go vue react docker}
}
