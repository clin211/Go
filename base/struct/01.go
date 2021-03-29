package main
// 结构体的基本使用方法
import (
	"fmt"
)

func main() {
	type Person struct {
		Nickname string
		Age      int
		Sign     string
		Gender   int
		Avatar   string
		Mobile   int
		Email    string
		location string
	}

	var p Person
	p.Nickname = "Forest"
	p.Age = 22
	p.Sign = "君子不器"
	p.Gender = 0
	p.Avatar = ""
	p.Mobile = 18782735415
	p.Email = "767425412@qq.com"
	p.location = "四川省成都市高新区天府软件园"

	fmt.Printf("个人信息：昵称:%s 年龄:%d 签名:%s 地址:%s", p.Nickname, p.Age, p.Sign, p.location)
}
