package main

import "fmt"

func Concat(a, b interface{}) string {
	return fmt.Sprintf("%v%v", a, b)
}
