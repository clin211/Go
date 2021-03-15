package main

import (
	"fmt"
)

func main() {
	// combination()
	res := combination()
	fmt.Println("res:", res)
	// num := combinationReturnArray()
	// fmt.Println("共有", num, "种组合方式")
}

// 有 1、2、3、4 这四个数字，能组成多少个互不相同且无重复数字的三位数？都是多少？
func combination() []int {
	count := 0
	res := []int{}
	for i := 1; i < 5; i++ {
		for j := 1; j < 5; j++ {
			for k := 1; k < 5; k++ {
				if i != j && i != k && k != j {
					count += 1
					val := []int{i, j, k}
					res = append(res, val...)
				}
			}
		}
	}
	fmt.Println("共", count, "种组合方式")
	return res
}

// Todo 返回已经组合起来的二维数组
// func combinationReturnArray() int {
// 	count := 0
// 	res := [...][3]int{}
// 	for i := 1; i < 5; i++ {
// 		for j := 1; j < 5; j++ {
// 			for k := 1; k < 5; k++ {
// 				if i != j && i != k && k != j {
// 					res[count][0] = i
// 					res[count][1] = j
// 					res[count][2] = k
// 					count += 1	// count++
// 				}
// 			}
// 		}
// 	}
// 	fmt.Println("res:", res)
// 	return count
// }
