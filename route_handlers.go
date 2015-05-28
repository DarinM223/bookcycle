package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

/*
 * Template parameter types
 */

// UserTemplateType For displaying navigation bar and anything that depends on the logged in User
type UserTemplateType struct {
	CurrentUser    User
	HasCurrentUser bool
}

// MessageTemplateType is a struct for the message template
type MessageTemplateType struct {
	UserTemplateType

	UserID int
}

// BookTemplateType is for displaying a book
type BookTemplateType struct {
	UserTemplateType

	Book      Book
	UserID    int
	CanDelete bool
}

// ManyBookTemplateType is for displaying many books (reused for many different things like
// search book results, recent books, and your books)
type ManyBookTemplateType struct {
	UserTemplateType

	Books []Book
	Title string
}

// GenerateFullTemplate Returns complete template with navigation bar added and your user login template
func GenerateFullTemplate(r *http.Request, bodyTemplatePath string) (*template.Template, UserTemplateType, error) {
	currentUser, err := CurrentUser(r)
	hasCurrentUser := true
	if err != nil {
		hasCurrentUser = false
	}

	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", bodyTemplatePath)
	if err != nil {
		return nil, UserTemplateType{}, err
	}

	return t, UserTemplateType{
		currentUser,
		hasCurrentUser,
	}, nil
}

/*
 * Route Handlers
 */

// RootHandler Route: /
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
		var recentBooks []Book
		result := db.Order("created_at desc").Limit(10).Find(&recentBooks)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}

		t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, ManyBookTemplateType{
			UserTemplateType: params,
			Books:            recentBooks,
			Title:            "Recent books",
		})
	}
}

// LogoutHandler Route: /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := LogoutUser(r, w)
	if err != nil {
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// LoginHandler Route: /login
func LoginHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	r.ParseForm()
	emailField := r.PostFormValue("email")
	passwordField := r.PostFormValue("password")

	validateFn := func() (User, error) {
		var user User
		result := db.First(&user, "email = ?", emailField)
		if result.Error != nil {
			return User{}, errors.New("Email or password is incorrect")
		}

		if user.Validate(passwordField) {
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

// UserJSONHandler Route: /users/{id}/json
func UserJSONHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	userID := mux.Vars(r)["id"]
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

// MapSearchHandler Route: /map_search/{id}
func MapSearchHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	_, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	receiverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("templates/boilerplate/nothing_boilerplate.html", "templates/map_search.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, struct {
		ReceiverID int
	}{
		ReceiverID: receiverID,
	})
}
