package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Println("listening on", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
