package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"database/sql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// DBInject injects a database object into a http handler with the database object parameter and
// turns it into a standard http handler
func DBInject(fn func(http.ResponseWriter, *http.Request, gorm.DB), db gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

// Routes returns a router that includes all of the routes needed for the application
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
	r.Methods("GET").Path("/users/{id}/json").HandlerFunc(DBInject(UserJSONHandler, db))
	r.Methods("GET", "POST").Path("/books/new").HandlerFunc(DBInject(NewBookHandler, db))
	r.Methods("GET").Path("/books").HandlerFunc(DBInject(ShowBooksHandler, db))
	r.Methods("GET").Path("/books/{id}/delete").HandlerFunc(DBInject(DeleteBookHandler, db))
	r.Methods("GET").Path("/books/{id}").HandlerFunc(DBInject(BookHandler, db))
	r.Methods("GET").Path("/search_results.json").HandlerFunc(DBInject(SearchResultsJSONHandler, db))
	r.Methods("GET").Path("/course_search.json").HandlerFunc(DBInject(CourseSearchHandler, db))
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

// SeedCourses seeds the main database with the courses inside the seed database
func SeedCourses(mainDB gorm.DB, seedDB *sql.DB) error {
	var (
		id                 int
		department         string
		courseID           string
		professorLastName  string
		professorFirstName string
		createdAt          time.Time
		updatedAt          time.Time
	)

	rows, err := seedDB.Query("SELECT * FROM courses")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &department, &courseID,
			&professorLastName, &professorFirstName, &createdAt, &updatedAt)
		if err != nil {
			return err
		}

		course := Course{
			ID:                 id,
			Department:         department,
			CourseID:           courseID,
			ProfessorLastName:  professorLastName,
			ProfessorFirstName: professorFirstName,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}
		mainDB.Create(&course)
	}
	return nil
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
	db.AutoMigrate(&User{}, &Book{}, &Message{}, &Course{})

	if len(os.Args) > 1 {
		seed := os.Args[1]
		if seed == "seed" {
			fmt.Println("Seeding courses from course sqlite file:")
			courseDB, err := sql.Open("sqlite3", "./CS188")
			if err != nil {
				fmt.Println(err.Error())
			}

			err = SeedCourses(db, courseDB)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Finished seeding courses")
			return
		}
	} else {
		fmt.Println("Listening...")
		PORT := os.Getenv("PORT")
		if PORT == "" {
			PORT = "8080"
			os.Setenv("PORT", PORT)
		}
		http.ListenAndServe(":"+PORT, Routes(db))
	}
}
