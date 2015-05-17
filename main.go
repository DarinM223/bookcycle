package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// Injects a database object into a http handler with the database object parameter and
// turns it into a standard http handler
func DBInject(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func Routes(db gorm.DB) *mux.Router {
	// Set up login sessions
	InitSessions("bookcycle")

	// run websocket hub and set websocket handler to /ws route
	go h.run(db)

	// Define routes (route handlers are in route_handlers.go)
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWs)
	r.HandleFunc("/", DBInject(RootHandler, db))
	r.Methods("POST").Path("/login").HandlerFunc(DBInject(LoginHandler, db))
	r.Methods("GET").Path("/logout").HandlerFunc(LogoutHandler)
	r.Methods("GET", "POST").Path("/users/new").HandlerFunc(DBInject(NewUserNewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/users/edit").HandlerFunc(DBInject(NewUserEditTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}").HandlerFunc(DBInject(NewUserViewTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}/json").HandlerFunc(DBInject(UserJsonHandler, db))
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(DBInject(NewBookHandler, db))
	r.Methods("GET").Path("/books").HandlerFunc(DBInject(ShowBooksHandler, db))
	r.Methods("GET").Path("/books/{id}/delete").HandlerFunc(DBInject(DeleteBookHandler, db))
	r.Methods("GET").Path("/books/{id}").HandlerFunc(DBInject(BookHandler, db))
	r.Methods("GET").Path("/search_results.json").HandlerFunc(DBInject(SearchResultsJsonHandler, db))
	r.Methods("GET").Path("/search_results").HandlerFunc(DBInject(SearchResultsHandler, db))
	r.Methods("GET").Path("/unread_messages").HandlerFunc(DBInject(UnreadMessagesHandler, db))
	r.Methods("GET").Path("/past_messages/{id}").HandlerFunc(DBInject(PastMessagesHandler, db))
	r.Methods("GET").Path("/message/{id}").HandlerFunc(DBInject(ChatHandler, db))

	// Set up static images
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/").Handler(fs)

	return r
}

func main() {
	// Set up database
	db, err := gorm.Open("sqlite3", "./sqlite_file.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Book{}, &Message{})

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", Routes(db))
}
