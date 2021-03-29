package main
// 结构体赋值
import (
	"fmt"
)

func main() {
	type Person struct {
		Name string
		Age  int
		Sign string
	}

	// 列表式赋值
	var person = Person{
		"Forset",
		22,
		"君子不器",
	}

	// 键值对方式赋值
	var forest = Person{
		Name: "Forest",
		Age:  22,
		Sign: "君子不器",
	}

	// 结构体嵌套蛮实用键值对初始化
	// developer = struct {
	// 	p Person
	// 	Language string
	// 	skill string
	// }{
		
	// }

	fmt.Println("person:", person)
	fmt.Println("forest", forest)
}
