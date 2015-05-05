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

func DBInject(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func main() {
	InitSessions("bookcycle")

	db, err := gorm.Open("sqlite3", "./sqlite_file.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Book{})

	fs := http.FileServer(http.Dir("./static/"))
	r := mux.NewRouter()
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	r.HandleFunc("/", RootHandler)
	r.Methods("POST").Path("/login").HandlerFunc(DBInject(LoginHandler, db))
	r.Methods("GET").Path("/logout").HandlerFunc(LogoutHandler)
	r.Methods("GET", "POST").Path("/users/new").HandlerFunc(DBInject(NewUserNewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/users/{id}").HandlerFunc(DBInject(NewUserEditTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(DBInject(NewBookHandler, db))
	r.Methods("GET", "DELETE").Path("/books/{id}").HandlerFunc(DBInject(BookHandler, db))
	r.Methods("GET", "POST").Path("/search").HandlerFunc(SearchHandler)

	r.PathPrefix("/").Handler(fs)
	fmt.Println("Listening...")
	http.ListenAndServe(":8080", r)
}

func showLoginPage(w http.ResponseWriter) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, nil)
}

func showUserPage(w http.ResponseWriter, u User) {
	t, err := template.ParseFiles("templates/loggedin.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, u)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	user, err := CurrentUser(r, w)
	if err != nil {
		showLoginPage(w)
	} else {
		showUserPage(w, user)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := LogoutUser(r, w)
	if err != nil {
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "POST" {
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
	} else {
		http.NotFound(w, r)
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	book_id := mux.Vars(r)["id"]
	fmt.Println(book_id)
	if r.Method == "GET" {
		var book Book
		result := db.First(&book, book_id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}

		t, err := template.ParseFiles("templates/book_detail.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, book)
	} else if r.Method == "DELETE" {
		// TODO: implement this
	} else {
		http.NotFound(w, r)
	}
}

func NewBookHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/new_book.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// TODO: implement this
		current_user, err := CurrentUser(r, w)
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
