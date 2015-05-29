package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// DBInject injects a database object into a http handler with the database object parameter and
// turns it into a standard http handler
func DBInject(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func SetupDB(dbType string, dbFilePath string) (gorm.DB, error) {
	db, err := gorm.Open(dbType, dbFilePath)
	if err != nil {
		return gorm.DB{}, err
	}
	db.AutoMigrate(&User{}, &Book{}, &Message{}, &Course{})
	return db, nil
}

func init() {
	// Set up login sessions
	InitSessions("bookcycle")
}

// Routes returns a router that includes all of the routes needed for the application
func Routes(db gorm.DB) *mux.Router {
	InitSessions("bookcycle")

	// run websocket hub and set websocket handler to /ws route
	go h.run(db)

	// Define routes (route handlers are in route_handlers.go)
	r := mux.NewRouter()
	r.HandleFunc("/ws", ServeWs)
	r.HandleFunc("/", DBInject(RootHandler, db))
	r.Methods("POST").Path("/login").HandlerFunc(DBInject(LoginHandler, db))
	r.Methods("GET").Path("/logout").HandlerFunc(LogoutHandler)
	r.Methods("GET", "POST").Path("/users/new").HandlerFunc(DBInject(NewUserNewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/users/edit").HandlerFunc(DBInject(NewUserEditTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}").HandlerFunc(DBInject(NewUserViewTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}/json").HandlerFunc(DBInject(UserJSONHandler, db))
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(DBInject(NewBookHandler, db))
	r.Methods("GET").Path("/books").HandlerFunc(DBInject(ShowBooksHandler, db))
	r.Methods("GET").Path("/books/{id}/delete").HandlerFunc(DBInject(DeleteBookHandler, db))
	r.Methods("GET").Path("/books/{id}").HandlerFunc(DBInject(BookHandler, db))
	r.Methods("GET").Path("/search_results.json").HandlerFunc(DBInject(SearchResultsJSONHandler, db))
	r.Methods("GET").Path("/courses/{id}/json").HandlerFunc(DBInject(CoursesJSONHandler, db))
	r.Methods("GET").Path("/course_search.json").HandlerFunc(DBInject(CourseSearchHandler, db))
	r.Methods("GET").Path("/search_results").HandlerFunc(DBInject(SearchResultsHandler, db))
	r.Methods("GET").Path("/messages").HandlerFunc(DBInject(MessagesHandler, db))
	r.Methods("GET").Path("/past_messages/{id}").HandlerFunc(DBInject(PastMessagesHandler, db))
	r.Methods("GET").Path("/message/{id}").HandlerFunc(DBInject(ChatHandler, db))
	r.Methods("GET").Path("/map_search/{id}").HandlerFunc(DBInject(MapSearchHandler, db))

	// Set up static images
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/").Handler(fs)

	return r
}
