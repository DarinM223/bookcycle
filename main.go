package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	// ./static/css/main.css maps to
	// localhost:blah/public/css/main.css
	http.Handle("/public", fs)
	http.HandleFunc("/", RootHandler)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, nil)
}
