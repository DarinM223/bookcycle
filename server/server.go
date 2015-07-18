package server

import (
	"github.com/PuerkitoBio/throttled"
	throttledStore "github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/justinas/nosurf"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

type DBInjectFunc func(func(http.ResponseWriter, *http.Request, gorm.DB), gorm.DB) http.Handler

func DBInject(requestsPerMinute int, testing bool) DBInjectFunc {
	st := throttledStore.NewMemStore(1000)
	return func(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) http.Handler {
		if testing {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fn(w, r, db)
			})
		}
		t := throttled.RateLimit(throttled.Q{Requests: requestsPerMinute, Window: time.Minute},
			&throttled.VaryBy{Path: true}, st)
		return t.Throttle(nosurf.NewPure(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, db)
		})))
	}
}

// DBInject injects a database object into a http handler with the database object parameter and
// turns it into a standard http handler

func init() {
	// Set up login sessions
	InitSessions("bookcycle")
}

// Routes returns a router that includes all of the routes needed for the application
func Routes(db gorm.DB, courseDB gorm.DB, requestsPerMinute int, testing bool) *mux.Router {
	InitSessions("bookcycle")

	// run websocket hub and set websocket handler to /ws route
	go h.run(db)

	DBInject := DBInject(requestsPerMinute, testing)

	// Define routes (route handlers are in route_handlers.go)
	r := mux.NewRouter()
	r.HandleFunc("/ws", ServeWs)
	r.Handle("/", DBInject(RootHandler, db))
	r.Methods("POST").Path("/login").Handler(DBInject(LoginHandler, db))
	r.Methods("GET").Path("/logout").HandlerFunc(LogoutHandler)
	r.Methods("GET", "POST").Path("/users/new").Handler(DBInject(NewUserNewTemplate().Handler, db))
	r.Methods("GET", "POST").Path("/users/edit").Handler(DBInject(NewUserEditTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}").Handler(DBInject(NewUserViewTemplate().Handler, db))
	r.Methods("GET").Path("/users/{id}/json").Handler(DBInject(UserJSONHandler, db))
	r.Methods("GET", "POST").Path("/books/new").Handler(DBInject(NewBookHandler, db))
	r.Methods("GET").Path("/books").Handler(DBInject(ShowBooksHandler, db))
	r.Methods("GET").Path("/books/{id}/delete").Handler(DBInject(DeleteBookHandler, db))
	r.Methods("GET").Path("/books/{id}").Handler(DBInject(BookHandler, db))
	r.Methods("GET").Path("/search_results.json").Handler(DBInject(SearchResultsJSONHandler, db))
	r.Methods("GET").Path("/courses/{id}/json").Handler(DBInject(CoursesJSONHandler, courseDB))
	r.Methods("GET").Path("/course_search.json").Handler(DBInject(CourseSearchHandler, courseDB))
	r.Methods("GET").Path("/search_results").Handler(DBInject(SearchResultsHandler, db))
	r.Methods("GET").Path("/messages").Handler(DBInject(MessagesHandler, db))
	r.Methods("GET").Path("/past_messages/{id}").Handler(DBInject(PastMessagesHandler, db))
	r.Methods("GET").Path("/message/{id}").Handler(DBInject(ChatHandler, db))
	r.Methods("GET").Path("/map_search/{id}").Handler(DBInject(MapSearchHandler, db))

	// Set up static images
	// ./static/css/main.css maps to
	// localhost:blah/css/main.css
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/").Handler(fs)

	return r
}
