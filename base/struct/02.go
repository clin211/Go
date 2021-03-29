// 结构体嵌套案例
package main

import (
	"fmt"
)

func main() {
	type Person struct {
		Name     string
		Age      int
		Gender   int
		Location string
	}

	type Student struct {
		P        Person
		Classics string
		Teacher  string
		Score    float64
	}

	var stu Student
	stu.P.Name = "Forest"
	stu.P.Age = 22
	stu.P.Gender = 0
	stu.P.Location = "成都"
	stu.Classics = "软件3班"
	stu.Teacher = "邱老师"
	stu.Score = 96.7

	fmt.Println("stu:", stu)

	//定义一个结构体嵌套，使用键值对初始化
	developer := Student{
		P: Person{
			Name: "Forest",
			Age:  20,
		},
		// Score: 60,
	}

	fmt.Println("Student:", developer) // developer {{Forest 20 0 }   60}
}
