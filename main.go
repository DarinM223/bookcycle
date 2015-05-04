package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

func main() {
	db, err := gorm.Open("sqlite3", "./sqlite_file.db")
	_ = db
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fs := http.FileServer(http.Dir("./static"))
	r := mux.NewRouter()
	// ./static/css/main.css maps to
	// localhost:blah/public/css/main.css
	http.Handle("/public/", fs)
	r.HandleFunc("/", RootHandler)
	r.Methods("GET", "DELETE").Path("/books/{id}").HandlerFunc(BookHandler)

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

func BookHandler(w http.ResponseWriter, r *http.Request) {
	book_id := mux.Vars(r)["id"]
	fmt.Println(book_id)
	test_book := Book{
		Title:     "Sample book",
		Author:    "Sample author",
		Class:     "Sample class",
		Professor: "Sample professor",
		Version:   "Sample version",
		Price:     "Sample price",
		Condition: "Sample condition",
	}

	t, err := template.ParseFiles("templates/book_detail.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, test_book)
}
