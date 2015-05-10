package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
)

/*
 * Template parameter types
 */

// For displaying navigation bar and anything that depends on the logged in User
type UserTemplateType struct {
	CurrentUser    User
	HasCurrentUser bool
}

type MessageTemplateType struct {
	UserTemplateType

	UserId int
}

// For displaying a book
type BookTemplateType struct {
	UserTemplateType

	Book      Book
	UserId    int
	CanDelete bool
}

// For displaying many books (reused for many different things like
// search book results, recent books, and your books)
type ManyBookTemplateType struct {
	UserTemplateType

	Books []Book
	Title string
}

// Returns complete template with navigation bar added and your user login template
func GenerateFullTemplate(r *http.Request, bodyTemplatePath string) (*template.Template, UserTemplateType, error) {
	current_user, err := CurrentUser(r)
	has_current_user := true
	if err != nil {
		has_current_user = false
	}

	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", bodyTemplatePath)
	if err != nil {
		return nil, UserTemplateType{}, err
	}

	return t, UserTemplateType{
		current_user,
		has_current_user,
	}, nil
}

/*
 * Route Handlers
 */

// Route: /
func RootHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	_, err := CurrentUser(r)
	if err != nil { // show login page if not logged in
		t, err := template.ParseFiles("templates/boilerplate/normal_boilerplate.html", "templates/index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		t.Execute(w, nil)
	} else { // show recent book listings if logged in
		var recent_books []Book
		result := db.Order("created_at desc").Limit(10).Find(&recent_books)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}

		t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
		if err != nil {
			return
		}

		t.Execute(w, ManyBookTemplateType{
			UserTemplateType: params,
			Books:            recent_books,
			Title:            "Recent books",
		})
	}
}

// Route: /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := LogoutUser(r, w)
	if err != nil {
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Route: /login
func LoginHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	r.ParseForm()
	email_field := r.PostFormValue("email")
	password_field := r.PostFormValue("password")

	validateFn := func() (User, error) {
		var user User
		result := db.First(&user, "email = ?", email_field)
		if result.Error != nil {
			return User{}, errors.New("Email or password is incorrect")
		}

		if user.Validate(password_field) {
			return user, nil
		}
		return User{}, errors.New("Email or password is incorrect")
	}

	err := LoginUser(r, w, validateFn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Route: /books
func ShowBooksHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var my_books []Book
	result := db.Model(&params.CurrentUser).Related(&my_books)
	if result.Error != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
	}

	t.Execute(w, ManyBookTemplateType{
		UserTemplateType: params,
		Books:            my_books,
		Title:            "My books",
	})
}

// Route: /books/{id}
func BookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	book_id := mux.Vars(r)["id"]
	var book Book
	result := db.First(&book, book_id)
	if result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	t, params, err := GenerateFullTemplate(r, "templates/book_detail.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	can_delete := params.CurrentUser.Id == book.UserId

	t.Execute(w, BookTemplateType{
		UserTemplateType: params,
		Book:             book,
		UserId:           book.UserId,
		CanDelete:        can_delete,
	})
}

// Route: /books/{id}/delete
func DeleteBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	book_id := mux.Vars(r)["id"]
	current_user, err := CurrentUser(r)
	if err != nil {
		http.Error(w, "You have to logged in to delete books", http.StatusUnauthorized)
		return
	}

	var book Book
	result := db.First(&book, book_id)
	if result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	if book.UserId != current_user.Id {
		http.Error(w, "You cannot delete books that you do not own", http.StatusUnauthorized)
		return
	}

	db.Delete(&book)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Route: /books/new
func NewBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		t, params, err := GenerateFullTemplate(r, "templates/new_book.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		t.Execute(w, params)
	} else if r.Method == "POST" {
		current_user, err := CurrentUser(r)
		if err != nil {
			http.Error(w, "You have to be logged in to add a book", http.StatusUnauthorized)
			return
		}

		book, err := NewMuxBookFactory().NewFormBook(r, current_user.Id)
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

// Route /search_results?query=
func SearchResultsHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	var search_books []Book
	result := db.Where("title LIKE ?", "%"+query+"%").Limit(10).Find(&search_books)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, ManyBookTemplateType{
		UserTemplateType: params,
		Books:            search_books,
		Title:            "Search Results",
	})
}

// Route: /unread_messages
func UnreadMessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	current_user, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var recentMessages []Message
	db.Where("receiver_id = ? and read = ?", current_user.Id, false).
		Order("created_at desc").Limit(10).Find(&recentMessages)
	if len(recentMessages) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
		return
	}
	messages_json, err := json.Marshal(recentMessages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(messages_json)
}

// Route: /message/{id}
func ChatHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	receiver_id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
	}

	t, params, err := GenerateFullTemplate(r, "templates/chat.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, MessageTemplateType{
		UserTemplateType: params,
		UserId:           receiver_id,
	})
}
