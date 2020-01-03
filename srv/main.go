package main

import (
	"waitqueue/srv/handler"
)
func main() {
	//r := mux.NewRouter()
	// mux 的经典构造方法具体使用
	// https://github.com/gorilla/mux/edit/master/README.md
	// Registered URLs
	/*
	r.Host(*addr).
		Path("/Login/{userId:[0-9]+}").
		Queries("userId", "{userId}").
		HandlerFunc(handler.Login).
		Name("Login")
	*/
	/*
	r.HandleFunc("/Login/{userId:[0-9]+}", handler.Login).
		Name("Login")
	log.Fatal(http.ListenAndServe(*addr, r))
	*/
	done := make(chan int)
	go handler.RunLoginServer()
	go handler.RunTokenServer()
	<-done
}
