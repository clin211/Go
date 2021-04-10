package main

import "fmt"

func main() {
	var year, months, days int
	fmt.Print("请输入你的出生年月：")
	// 获取输入的年月日信息
	fmt.Scanf("%d年%d月%d日\n", &year, &months, &days)
	// 打印输入的年月日
	fmt.Println("birthday:", year)

	fmt.Print("数据拆分：", year, months, days)
}

