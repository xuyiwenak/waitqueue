package main

import (
	"flag"
	"log"
	"net/http"
	"waitqueue/srv/handler"
)

var addr = flag.String("addr", "172.16.21.23:8083", "http service address")


func main() {
	flag.Parse()
	log.SetFlags(0)
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
	go handler.RunTokenServer()
	handler.InitSeqQueue(10000)
	http.HandleFunc("/Login", handler.Login)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
