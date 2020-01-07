package main

import (
	"waitqueue/srv/handler"
)
func main() {
	done := make(chan int)
	go handler.RunLoginServer()
	go handler.RunTokenServer()
	<-done
}
