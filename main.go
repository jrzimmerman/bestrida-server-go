package main

import (
	"log"
	"net/http"
)

const addr = ":4001"

func main() {
	log.Println("listening on", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
