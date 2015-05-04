package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	r := mux.NewRouter()
	// ./static/css/main.css maps to
	// localhost:blah/public/css/main.css
	r.Handle("/public", fs)
	r.HandleFunc("/", RootHandler)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", r)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, nil)
}
