package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/about.html")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

