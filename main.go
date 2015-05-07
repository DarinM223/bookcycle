package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

// Injects a database object into a http handler with the database object parameter and
// turns it into a standard http handler
func DBInject(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func main() {
	// Set up login sessions
	InitSessions("bookcycle")

	// Set up database
	db, err := gorm.Open("sqlite3", "./sqlite_file.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Book{})

	// Define routes
	r := mux.NewRouter()
	r.HandleFunc("/", DBInject(RootHandler, db))
	r.Methods("POST").Path("/login").HandlerFunc(DBInject(LoginHandler, db))
	r.Methods("GET").Path("/logout").HandlerFunc(LogoutHandler)
	r.Methods("GET", "POST").Path("/users/new").HandlerFunc(DBInject(NewUserNewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/users/edit").HandlerFunc(DBInject(NewUserEditTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}").HandlerFunc(DBInject(NewUserViewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(DBInject(NewBookHandler, db))
	r.Methods("GET").Path("/books/{id}/delete").HandlerFunc(DBInject(DeleteBookHandler, db))
	r.Methods("GET").Path("/books/{id}").HandlerFunc(DBInject(BookHandler, db))
	r.Methods("GET").Path("/search").HandlerFunc(SearchHandler)
	r.Methods("GET").Path("/search_results").HandlerFunc(DBInject(SearchResultsHandler, db))

	// Set up static images
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/").Handler(fs)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", r)
}

func showLoginPage(w http.ResponseWriter) {
	t, err := template.ParseFiles("templates/boilerplate/normal_boilerplate.html", "templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, nil)
}

func showUserPage(w http.ResponseWriter, u User, db gorm.DB) {
	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", "templates/search_results.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var recent_books []Book
	result := db.Order("created_at desc").Limit(10).Find(&recent_books)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	t.Execute(w, struct {
		CurrentUser    User
		HasCurrentUser bool
		RecentBooks    []Book
		Title          string
	}{
		CurrentUser:    u,
		HasCurrentUser: true,
		RecentBooks:    recent_books,
		Title:          "Recent Books",
	})
}

// Route: /
func RootHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	user, err := CurrentUser(r)
	if err != nil {
		showLoginPage(w)
	} else {
		showUserPage(w, user, db)
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

// Route: /books/{id}
func BookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	current_user, err := CurrentUser(r)
	has_current_user := true
	if err != nil {
		has_current_user = false
	}
	book_id := mux.Vars(r)["id"]
	var book Book
	result := db.First(&book, book_id)
	if result.Error != nil {
		http.Error(w, "Book does not exist", http.StatusUnauthorized)
		return
	}

	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", "templates/book_detail.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	can_delete := current_user.Id == book.UserId
	t.Execute(w, struct {
		Book           Book
		UserId         int
		CurrentUser    User
		HasCurrentUser bool
		CanDelete      bool
	}{
		Book:           book,
		UserId:         book.UserId,
		CurrentUser:    current_user,
		HasCurrentUser: has_current_user,
		CanDelete:      can_delete,
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
		current_user, err := CurrentUser(r)
		has_current_user := true
		if err != nil {
			has_current_user = false
		}
		t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
			"templates/navbar.html", "templates/new_book.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, struct {
			CurrentUser    User
			HasCurrentUser bool
		}{
			current_user,
			has_current_user,
		})
	} else if r.Method == "POST" {
		current_user, err := CurrentUser(r)
		if err != nil {
			http.Error(w, "You have to be logged in to add a book", http.StatusUnauthorized)
			return
		}
		book, err := NewMuxBookFactory().NewFormBook(r, current_user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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

// Route: /search
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	current_user, err := CurrentUser(r)
	has_current_user := true
	if err != nil {
		has_current_user = false
	}
	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", "templates/search.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, struct {
		CurrentUser    User
		HasCurrentUser bool
	}{
		current_user,
		has_current_user,
	})
}

// Route /search_results?query=
func SearchResultsHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}
	current_user, err := CurrentUser(r)
	has_current_user := true
	if err != nil {
		has_current_user = false
	}
	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", "templates/search_results.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var search_books []Book
	result := db.Where("title LIKE ?", "%"+query+"%").Limit(10).Find(&search_books)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	t.Execute(w, struct {
		CurrentUser    User
		HasCurrentUser bool
		RecentBooks    []Book
		Title          string
	}{
		current_user,
		has_current_user,
		search_books,
		"Search Results",
	})
}
