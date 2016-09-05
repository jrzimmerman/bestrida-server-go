package routes

import (
	"fmt"
	"net/http"
)

// HelloHandler returns "Hello World" to the request
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func routes() {
	http.HandleFunc("/", HelloHandler)
}
