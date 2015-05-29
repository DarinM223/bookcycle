package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ShowBooksHandler Route: /books
func ShowBooksHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var myBooks []Book
	result := db.Model(&params.CurrentUser).Related(&myBooks)
	if result.Error != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}

	t.Execute(w, ManyBookTemplateType{
		UserTemplateType: params,
		Books:            myBooks,
		Title:            "My books",
	})
}

// BookHandler Route: /books/{id}
func BookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	bookID := mux.Vars(r)["id"]
	var book Book
	result := db.First(&book, bookID)
	if result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	t, params, err := GenerateFullTemplate(r, "templates/book_detail.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	canDelete := params.CurrentUser.ID == book.UserID

	t.Execute(w, BookTemplateType{
		UserTemplateType: params,
		Book:             book,
		UserID:           book.UserID,
		CanDelete:        canDelete,
	})
}

// DeleteBookHandler Route: /books/{id}/delete
func DeleteBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	bookID := mux.Vars(r)["id"]
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.Error(w, "You have to logged in to delete books", http.StatusUnauthorized)
		return
	}

	var book Book
	result := db.First(&book, bookID)
	if result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	if book.UserID != currentUser.ID {
		http.Error(w, "You cannot delete books that you do not own", http.StatusUnauthorized)
		return
	}

	db.Delete(&book)
	http.Redirect(w, r, "/", http.StatusFound)
}

type isbnSearchResult struct {
	TotalItems int              `json:"totalItems"`
	Items      []isbnBookResult `json:"items"`
}

type isbnBookResult struct {
	ID         string         `json:"id"`
	Etag       string         `json:"etag"`
	VolumnInfo isbnVolumeInfo `json:"volumnInfo"`
}

type isbnVolumeInfo struct {
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Publisher     string   `json:"publisher"`
	PublishedDate string   `json:"publishedDate"`
	Description   string   `json:"description"`
	AverageRating float64  `json:"averageRating"`
}

// NewBookISBNHandler Route: POST /books/new.json?isbn=
func NewBookISBNHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	isbn := r.URL.Query().Get("isbn")
	if len(isbn) == 0 {
		http.Error(w, "You must include an ISBN parameter", http.StatusNotFound)
		return
	}

	// Send request to google book api with isbn
	// Example request: https://www.googleapis.com/books/v1/volumes?q=isbn:0735619670
	res, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=isbn:" + url.QueryEscape(isbn))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// read response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var searchResult isbnSearchResult
	err = json.Unmarshal(data, &searchResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Grab relevant data and store it into database
}

// NewBookHandler Route: /books/new
func NewBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		t, params, err := GenerateFullTemplate(r, "templates/new_book.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		t.Execute(w, params)
	} else if r.Method == "POST" {
		currentUser, err := CurrentUser(r)
		if err != nil {
			http.Error(w, "You have to be logged in to add a book", http.StatusUnauthorized)
			return
		}

		book, err := NewMuxBookFactory().NewFormBook(r, currentUser.ID)
		if err != nil {
			http.Error(w, "There was an error with validating some of your fields. Please check your input again",
				http.StatusUnauthorized)
			return
		}

		result := db.Create(&book)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}
