package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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

var store = sessions.NewCookieStore([]byte("helloworld"))

func main() {
	store.Options = &sessions.Options{
		//Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}
	gob.Register(&User{})

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
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(NewBookHandler)
	r.Methods("GET", "DELETE").Path("/books/{id}").HandlerFunc(BookHandler)
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
	sess, err := store.Get(r, "bookcycle")
	if err != nil {
		showLoginPage(w)
	} else {
		if user, ok := sess.Values["user"]; ok {
			if user != nil {
				showUserPage(w, *user.(*User))
			} else {
				showLoginPage(w)
			}
		} else {
			showLoginPage(w)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "bookcycle")
	if err != nil {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
	} else {
		if _, ok := sess.Values["user"]; ok {
			delete(sess.Values, "user")
			err := sess.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		email_field := r.PostFormValue("email")
		password_field := r.PostFormValue("password")

		sess, err := store.Get(r, "bookcycle")
		if err != nil {
			sess, err = store.New(r, "bookcycle")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		}
		if _, ok := sess.Values["user"]; ok {
			http.Error(w, "User is already logged in!", http.StatusUnauthorized)
		} else {
			var user User
			result := db.First(&user, "email = ?", email_field)
			if result.Error != nil {
				http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
			} else {
				if user.Validate(password_field) {
					sess.Values["user"] = user
					err := sess.Save(r, w)
					if err != nil {
						http.Error(w, err.Error(), http.StatusUnauthorized)
						return
					}
					http.Redirect(w, r, "/", http.StatusFound)
				} else {
					http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
				}
			}
		}
	} else {
		http.NotFound(w, r)
	}
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
