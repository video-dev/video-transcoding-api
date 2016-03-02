package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./swagger-ui-server")))
	http.ListenAndServe(":7777", nil)
}
