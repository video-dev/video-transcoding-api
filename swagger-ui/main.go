package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./swagger-ui")))
	http.ListenAndServe(":8888", nil)
}
