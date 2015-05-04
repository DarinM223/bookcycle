package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	http.Handle("/", fs)
	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
