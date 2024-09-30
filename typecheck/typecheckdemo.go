package typecheck

import "fmt"

func check() {
	t := unexistType{}
	fmt.Println(t)
}

func unused() {
	i := 1
}
