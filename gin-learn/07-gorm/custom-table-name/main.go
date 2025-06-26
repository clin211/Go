package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
}

func (User) TableName() string {
	return "users"
}

func main() {
	db, err := gorm.Open(mysql.Open("gorm:gorm123456@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Successfully connected to database!")

	db.AutoMigrate(&User{})

	db.Create(&User{Name: "John", Age: 20})

	var user User
	db.First(&user, 1)
	fmt.Println(user) // {{1 2025-05-19 20:25:23.411 +0800 CST 2025-05-19 20:25:23.411 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} John 20}
	json, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Failed to marshal user:", err)
	}
	fmt.Println(string(json)) // {"ID":1,"CreatedAt":"2025-05-19T20:25:23.411+08:00","UpdatedAt":"2025-05-19T20:25:23.411+08:00","DeletedAt":null,"Name":"John","Age":20}
}
