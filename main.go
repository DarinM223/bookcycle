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

	fs := http.FileServer(http.Dir("./static/"))
	r := mux.NewRouter()
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	r.HandleFunc("/", RootHandler)
	r.Methods("GET", "POST").Path("/users/new").HandlerFunc(NewUserNewTemplate().Handler)
	r.Methods("GET", "POST").Path("/users/{id}").HandlerFunc(NewUserEditTemplate().Handler)
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(NewBookHandler)
	r.Methods("GET", "DELETE").Path("/books/{id}").HandlerFunc(BookHandler)
	r.Methods("GET", "POST").Path("/search").HandlerFunc(SearchHandler)

	r.PathPrefix("/").Handler(fs)
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
	if r.Method == "GET" {
		// Test book to test that it populates the template with fields
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
	} else if r.Method == "DELETE" {
		// TODO: implement this
	} else {
		http.NotFound(w, r)
	}
}

func NewBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/new_book.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// TODO: implement this
	} else {
		http.NotFound(w, r)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // display search page
		t, err := template.ParseFiles("templates/search.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" { // display search results page
		t, err := template.ParseFiles("templates/search_results.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// TODO: set the template variable so you can show some actual search results
		t.Execute(w, nil)
	} else {
		http.NotFound(w, r)
	}
}
