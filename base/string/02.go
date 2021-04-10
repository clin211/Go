package main

import "fmt"
// 使用 Golang 输入爱心图案，爱心的公式为 （x²+y²-1）³-x²*y³=0。
func main() {
	// (x^2 + y^2 - 1)^3 = x^2y^3
	for y := 1.5; y > -1.5; y -= 0.1 {
		for x := -1.5; x < 1.5; x += 0.04 {
			a := x*x + y*y - 1
			if a*a*a-x*x*y*y*y <= 0.0 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}
