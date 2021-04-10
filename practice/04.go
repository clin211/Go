package main

// 输入某年某月某日，判断这一天是这一年的第几天？
import (
	"fmt"
)

func main() {
	// var y, m, d int = 0, 0, 0
	// var days int = 0
	// fmt.Printf("请输入年月日\n")
	// fmt.Scanf("%d%d%d\n", &y, &m, &d)
	// fmt.Printf("%d 年 %d 月 %d 日\n", y, m, d)
	// fmt.Println("输入的年月;", y, m, d)
	// day := [7]int{1, 3, 5, 7, 8, 10, 12}
	// fmt.Println("day:", day)
	// switch m {
	// case 12:
	// 	days += d
	// 	d = 30
	// 	fallthrough
	// case 11:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 10:
	// 	days += d
	// 	d = 30
	// 	fallthrough
	// case 9:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 8:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 7:
	// 	days += d
	// 	d = 30
	// 	fallthrough
	// case 6:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 5:
	// 	days += d
	// 	d = 30
	// 	fallthrough
	// case 4:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 3:
	// 	days += d
	// 	d = 28
	// 	if (y%400 == 0) || (y%4 == 0 && y%100 != 0) {
	// 		d += 1
	// 	}
	// 	fallthrough
	// case 2:
	// 	days += d
	// 	d = 31
	// 	fallthrough
	// case 1:
	// 	days += d
	// }
	// fmt.Printf("是今年的第 %d 天!\n", days)
	// calcDays()
	calcDaysArray()
}

func calcDays() {
	var year, month, day int
	fmt.Printf("请输入年月日；示例：2021年4月10日\n")
	fmt.Scanf("%d年%d月%d日\n", &year, &month, &day)
	monthly := [7]int{1, 3, 5, 7, 8, 10, 12}
	// 是否是闰年 闰年29天
	abortion := [4]int{4, 6, 9, 11}

	if year == 0 || month == 0 || day == 0 {
		return
	}
	var result int

	for _, v := range monthly {
		if v <= month {
			result += 31
		}
	}
	for _, v := range abortion {
		if v <= month {
			result += 30
		}
	}
	/*
	 * 满足以下两点中任意一点就是闰年
	 * - 能被4整除，但是不能被100整出。
	 * - 能被400整除
	 */

	if year&4 == 0 && year%100 != 0 || year%400 == 0 {
		fmt.Println("是闰年：", year)
		result += 29
	} else {
		result += 28
	}
	fmt.Println("第%d天", result)
}

func calcDaysArray() {
	var year, month, day int
	fmt.Printf("请输入年月日；示例：（2021年4月10日）\n")
	fmt.Scanf("%d年%d月%d日\n", &year, &month, &day)

	var months = [11][2]int{
		{1, 31},
		{3, 31},
		{4, 30},
		{5, 31},
		{6, 30},
		{7, 31},
		{8, 31},
		{9, 30},
		{10, 31},
		{11, 30},
		{12, 31},
	}
	var result int = day
	for i := 0; i < len(months); i++ {
		if months[i][0] < month {
			result += months[i][1]
			fmt.Println("result:", result)
		}
	}

	// year&4 == 0 && year%100 != 0 || year%400 == 0 ? result += 29 : result += 28
	if month > 2 {
		if year&4 == 0 && year%100 != 0 || year%400 == 0 {
			result += 29
		} else {
			result += 28
		}
	}

	fmt.Println("result:", result)
	fmt.Printf("%d年%d月%d日是%d的第%d天\n", year, month, day, result)
}
