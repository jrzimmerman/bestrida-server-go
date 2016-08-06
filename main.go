package main

import (
	"fmt"
	"log"
	"net/http"
)

const addr = ":4001"

// HelloHandler returns "Hello World" to the request
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func main() {
	http.HandleFunc("/", HelloHandler)
	log.Println("listening on", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
