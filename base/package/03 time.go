package main

import (
	. "fmt"
	"time"
)

func main() {
	// 格式化时间必须使用 2006 1 2 15:04:05 这几个数字
	t := time.Now()
	Println("当前时间：", t.Format("2006年1月2日 15:04:05"))

	// 按短横线连接
	s := t.Format("2006-1-2 15:04:05")
	Println(s) //打印出的格式就是当前的时间 2021-4-13 15:17:28

	// 用/连接
	s1 := t.Format("2006/1/2")
	Println(s1) //打印出的格式就是当前的年月日 2021/4/13

	// 字符串类型转换为时间类型
	str := "2021年4月15日"
	t, err := time.Parse("2006年1月2日", str)
	if err != nil {
		Println("err", err)
	}
	Println("时间转换", t)

	//获取年月日信息
	year, month, day := time.Now().Date()
	Println(year, month, day) // 2021 April 13

	//获取时分秒信息
	hour, minute, second := time.Now().Clock()
	Println(hour, minute, second) // 15 25 11

	//获取今年过了多少天了
	today := time.Now().YearDay()
	Println(today) // 103  (今年已经过了103天了)

	//获取今天是星期几
	weekday := time.Now().Weekday()
	Println(weekday) //Tuesday
}
