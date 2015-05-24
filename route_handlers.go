package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

// SearchCourse helper function for searching courses
func SearchCourse(department string, courseID string, professor string, db gorm.DB) ([]Course, error) {
	var searchCourses []Course
	if len(department) == 0 {
		return []Course{}, nil
	} else if len(courseID) == 0 {
		result := db.Where("department LIKE ?", "%"+department+"%").Limit(10).Find(&searchCourses)
		if result.Error != nil {
			return []Course{}, result.Error
		}
		return searchCourses, nil
	} else if len(professor) == 0 {
		result := db.Where("department LIKE ? AND course_id LIKE ?", "%"+department+"%", "%"+courseID+"%").
			Limit(10).Find(&searchCourses)
		if result.Error != nil {
			return []Course{}, result.Error
		}
		return searchCourses, nil
	} else {
		var result *gorm.DB
		professorArray := strings.Split(professor, " ")
		if len(professorArray) < 2 {
			result = db.Where(`department LIKE ? 
							AND course_id LIKE ? 
							AND professor_first_name LIKE ?`,
				"%"+department+"%", "%"+courseID+"%", "%"+professorArray[0]+"%").
				Limit(10).Find(&searchCourses)
		} else {
			result = db.Where(`department LIKE ? 
							AND course_id LIKE ? 
							AND professor_first_name LIKE ? 
							AND professor_last_name LIKE ?`,
				"%"+department+"%", "%"+courseID+"%", "%"+professorArray[0]+"%", "%"+professorArray[1]+"%").
				Limit(10).Find(&searchCourses)
		}
		if result.Error != nil {
			return []Course{}, result.Error
		}
		return searchCourses, nil
	}
	return []Course{}, nil
}

// SearchResultsJSONHandler Route /search_results.json?query=
func SearchResultsJSONHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	http.NotFound(w, r)

	//searchBooks, err := SearchBook(query, db)
	//if err != nil {
	//    http.NotFound(w, r)
	//    return
	//}

	//searchBooksJSON, err := json.Marshal(searchBooks)
	//if err != nil {
	//    http.NotFound(w, r)
	//    return
	//}

	//w.Header().Set("Content-Type", "application/json")
	//w.Write(searchBooksJSON)
}

// SearchResultsHandler Route /search_results?query=
func SearchResultsHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	// TODO: search google books api, get results, for each isbn check if book has seller

	http.NotFound(w, r)
}

func CourseSearchHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	department := r.URL.Query().Get("department")
	courseID := r.URL.Query().Get("course_id")
	professor := r.URL.Query().Get("professor")

	searchCourses, err := SearchCourse(department, courseID, professor, db)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	searchCoursesJSON, err := json.Marshal(searchCourses)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(searchCoursesJSON)
}

// UnreadMessagesHandler Route: /unread_messages
func UnreadMessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var recentMessages []Message
	db.Where("receiver_id = ? and read = ?", currentUser.ID, false).
		Order("created_at desc").Limit(10).Find(&recentMessages)
	if len(recentMessages) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
		return
	}
	messagesJSON, err := json.Marshal(recentMessages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(messagesJSON)
}

// PastMessagesHandler Route: /past_messages/{id}
func PastMessagesHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	receiverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var results []Message
	res := db.Where("(receiver_id = ? and sender_id = ?) or (receiver_id = ? and sender_id = ?)",
		currentUser.ID, receiverID, receiverID, currentUser.ID).Limit(20).Order("created_at desc").Find(&results)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resultsJSON)
}

// ChatHandler Route: /message/{id}
func ChatHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	receiverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if currentUser.ID == receiverID {
		http.Error(w, "You cannot message yourself", http.StatusUnauthorized)
		return
	}

	t, err := template.ParseFiles("templates/boilerplate/nothing_boilerplate.html",
		"templates/navbar.html", "templates/chat.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, MessageTemplateType{
		UserTemplateType: UserTemplateType{
			CurrentUser:    currentUser,
			HasCurrentUser: true,
		},
		UserID: receiverID,
	})
}
