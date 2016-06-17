package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// ShowBooksHandler is a route for /books that displays all books that you own
func ShowBooksHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var myBooks []Book
	if result := db.Model(&params.CurrentUser).Related(&myBooks); result.Error != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}

	t.Execute(w, ManyBookTemplateType{
		UserTemplateType: params,
		Books:            myBooks,
		Title:            "My books",
	})
}

// BookHandler is a route for /books/{id} that displays a book with a certain ID
func BookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	bookID := mux.Vars(r)["id"]
	var book Book
	if result := db.First(&book, bookID); result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	t, params, err := GenerateFullTemplate(r, "templates/book_detail.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, BookTemplateType{
		UserTemplateType: params,
		Book:             book,
		UserID:           book.UserID,
		CanDelete:        params.CurrentUser.ID == book.UserID,
	})
}

// DeleteBookHandler is a route for /books/{id}/delete that deletes a book with a certain ID
// You have to be logged in and you can only delete your own books
func DeleteBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	bookID := mux.Vars(r)["id"]
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.Error(w, "You have to logged in to delete books", http.StatusUnauthorized)
		return
	}

	var book Book
	if result := db.First(&book, bookID); result.Error != nil {
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

// NewBookHandler is a route for /books/new that creates a new book
// GET /books/new displays the new book page
// POST /books/new creates a new book from post parameters:
// isbn string
// title string
// course_id integer
// price float
// condition integer
// details string
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

		if result := db.Create(&book); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}
