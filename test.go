package main

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

func main() {
	var a int64
	atomic.StoreInt64(&a, 0)
	atomic.AddInt64(&a, 1)
	//atomic.AddInt64(&a, -1)

	fmt.Println(reflect.TypeOf(a))
	fmt.Println(atomic.LoadInt64(&a))
}
