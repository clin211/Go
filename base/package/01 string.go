package main

import (
	. "fmt"
	"strings"
)

func main() {
	// 检查是否包含指定值 返回布尔值
	var s1 string = "ok let's go"
	Println("检查是否包含指定值:", strings.Contains(s1, "go")) // 检查是否包含指定值: true

	// 检查是否包含指定字符串中任意一个字符 返回布尔值
	s2 := "study go!"
	Println("检查是否包含指定字符串中任意一个字符:", strings.ContainsAny(s2, "o")) // 检查是否包含指定字符串中任意一个字符: true

	// 检查指定字符出现的次数 返回数字
	s3 := "to add module requirements and"
	Println("检查指定字符出现的次数:", strings.Count(s3, "a")) // 检查指定字符出现的次数: 2

	// 检查文本的前缀是否是指定值 返回布尔值
	s4 := "creating new go.mod"
	Println("检查文本的前缀：", strings.HasPrefix(s4, "cs")) // 检查文本的前缀： false

	// 检查文本的前缀是否是指定值 返回布尔值
	s5 := "creating new go.mod"
	Println("检查文本的前缀：", strings.HasSuffix(s5, "mod")) // 检查文本的前缀： true

	// 检查指定字符在字符串中第一次出现的位置 若存在返回相应的下标 若不存在返回-1
	s6 := "completed with 3 local objects."
	Println("检查指定字符在字符串中第一次出现的位置：", strings.Index(s6, "o")) // 检查指定字符在字符串中第一次出现的位置： 9

	//查找字符中任意一个字符第一次出现在字符串中的位置
	s7 := "completed with 3 local objects."
	Println("查找字符中任意一个字符出现在字符串中的位置：", strings.IndexAny(s7, " "))

	/*
	 TODO strings.Index() 与 strings.IndexAny()的区别

	*/

	// 检查指定字符在字符串中最后一次出现的位置 若存在返回相应的下标 若不存在返回-1
	s8 := "completed with 3 local objects."
	Println("检查指定字符在字符串中最后一次出现的位置：", strings.LastIndex(s8, " ")) // 检查指定字符在字符串中最后一次出现的位置： 22

	// 字符串拼接
	s9 := "learn"
	s10 := " go"
	s11 := s9 + s10
	Println("加号拼接：", s11) // 加号拼接： learn go

	s12 := []string{"learn", "go"}
	s13 := strings.Join(s12, " ") // 用空格连接
	Println("Join拼接：", s13)       // Join拼接： learn go

	// 字符串切割
	s14 := "Enumerating objects: 7, done."
	s15 := strings.Split(s14, ":") // 返回一个数组
	Println("字符串切割：", s15[0])      // 字符串切割： Enumerating objects

	s16 := "Enumerating objects: 7, done."
	s17 := strings.Replace(s16, ":", "-", 1) // 最后一个参数是替换的次数，为-1的时候全部替换
	Println("字符串替换：", s17)                   // 字符串替换：Enumerating objects- 7, done.
}
