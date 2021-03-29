package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Person struct {
		Name   string
		Age    int
		Gender int
		Sign   string
	}
	var person Person

	// 将结构体转成json
	jsonStr := `{"Name" : "Forest", "Age" : 22, "Gender" : 0, "Sign" : "君子不器！"}`
	if err := json.Unmarshal([]byte(jsonStr), &person); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("person:", person) // person: {Forest 22 0 君子不器！}
}
