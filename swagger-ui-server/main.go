package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./swagger-ui-server")))
	log.Printf("Starting the swagger-ui-server at http://localhost:7777")
	http.ListenAndServe(":7777", nil)
}
