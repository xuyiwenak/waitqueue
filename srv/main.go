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
	http.HandleFunc("/Login", handler.Login)
	log.Fatal(http.ListenAndServe(*addr, nil))
}