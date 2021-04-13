package main

import (
	. "fmt"
	"strconv"
)

/**
Tasking:
	- string 转 int
	- string 转 bool
	- string 转 float
	- int 转 string
	- int 转 bool
	- float 转 string
	- float 转 int
	- bool 转 string
	- bool 转 int
*/
func main() {
	// string 转 int
	str := "100"
	n, err := strconv.ParseInt(str, 10, 64)
	// n, err := strconv.Atoi(str, 10, 64)
	if err != nil {
		Println("string转int：", err)
		return
	}
	Println("string转int：", n) // 100

	// 字符串转布尔
	str1 := "true"
	b, err := strconv.ParseBool(str1)
	if err != nil {
		Println("err:", err)
	}
	Println("字符串转布尔：", b) // true

	// 字符串转浮点数
	str2 := "6.6"
	f, err := strconv.ParseFloat(str2, 64)
	if err != nil {
		Println("err:", err)
	}
	Println("字符串转浮点数：", f) // 6.6

	// 数字转字符串
	number := 100
	str3 := strconv.Itoa(number)
	Println("数字转字符串：", str3) // "100"

	// 数字转布尔 有效值为0和1，
	var number1 int = 0
	cb, err := strconv.ParseBool(strconv.Itoa(number1))
	if err != nil {
		Println("err:", err)
	}
	Println("数字转布尔：", cb) // false

	// 布尔转字符串
	const isB = true
	sb := strconv.FormatBool(isB)
	Println("布尔转字符串:", sb) // "true"

	// float 类型转 int32
	var floatNumber float64 = 3.1415926
	var pi = int64(floatNumber)
	Println("float类型转int64:", pi) // 3
}
