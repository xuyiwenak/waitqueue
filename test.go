package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a uint64
	a = 64
	fmt.Println(reflect.TypeOf(a))
}
