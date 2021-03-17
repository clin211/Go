package main

import (
	"fmt"
)


func main() {
	var random []int = []int{432, 345, 345, 543, 234, 435, 34, 3, 67, 10}
	asc := bubbleSort(random, "asc")
	fmt.Println("asc:", asc)   // [3 10 34 67 234 345 345 432 435 543]
	desc := bubbleSort(random, "desc")
	fmt.Println("desc:", desc) // [543 435 432 345 345 234 67 34 10 3]

	a, b, c := returnParam()
	fmt.Println("param:", a, b, c)
}

// 冒泡排序
/**
 * @param array 排序的数据源
 * @param types 排序类型 types参数选项: asc、desc
*/
func bubbleSort(array []int, types string) []int {
	for i := 0; i < len(array); i++ {
		for j := i + 1; j < len(array); j++ {
			// 升序
			if types == "asc" && array[i] > array[j] {
				array[i], array[j] = array[j], array[i]
			}
			// 降序
			if types == "desc" && array[i] < array[j] {
				array[i], array[j] = array[j], array[i]
			}
		}
	}
	return array
}

func returnParam() (a, b, c int) {
	a, b, c = 111, 222, 333
	return
}
