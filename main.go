package main

import (
	"fmt"
	"github.com/DarinM223/bookcycle/server"
	"github.com/lib/pq"
	"net/http"
	"os"
	"time"

	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const RequestsPerMinute int = 30

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

		course := server.Course{
			ID:         id,
			Department: department,
			CourseID:   courseID,
			Professor:  professorFirstName + " " + professorLastName,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		}
		mainDB.Create(&course)
	}
	return nil
}

// IsTesting returns true if there are any command line arguments with the
// value "loadtest" and false otherwise. It is used as a parameter to server.Routes()
// so that rate limiting and csrf are turned off when "loadtest" is a command line argument
func IsTesting(args []string) bool {
	for _, arg := range args {
		if arg == "loadtest" {
			fmt.Println("Load testing enabled")
			return true
		}
	}
	return false
}

func main() {
	var db gorm.DB

	coursesDB, err := gorm.Open("sqlite3", "./courses.database")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer coursesDB.Close()
	coursesDB.AutoMigrate(&server.Course{})

	if len(os.Args) > 1 {
		option := os.Args[1]
		if option == "production" { // configure postgres database
			fmt.Println("Running in production mode")
			url := os.Getenv("DATABASE_URL")
			connection, _ := pq.ParseURL(url)
			connection += " sslmode=require"
			db, err = gorm.Open("postgres", connection)
			if err != nil {
				fmt.Println(err)
				return
			}
			db.AutoMigrate(&server.User{}, &server.Book{}, &server.Message{})
		} else if option == "seed" {
			fmt.Println("Seeding courses from course sqlite file:")
			db.LogMode(true)
			seedDB, err := sql.Open("sqlite3", "./CS188")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if err = SeedCourses(coursesDB, seedDB); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Finished seeding courses")
			return
		}
	} else { // configure sqlite database
		db, err = gorm.Open("sqlite3", "./sqlite_file.db")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		db.AutoMigrate(&server.User{}, &server.Book{}, &server.Message{})
	}
	fmt.Println("Listening...")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
		os.Setenv("PORT", PORT)
	}
	http.ListenAndServe(":"+PORT, server.Routes(db, coursesDB,
		RequestsPerMinute, IsTesting(os.Args)))
}
