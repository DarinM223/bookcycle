package main

import (
	"fmt"
	"github.com/DarinM223/bookcycle/server"
	"net/http"
	"os"
	"time"

	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

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

func main() {
	// Set up database
	db, err := gorm.Open("sqlite3", "./sqlite_file.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	db.AutoMigrate(&server.User{}, &server.Book{}, &server.Message{})

	coursesDB, err := gorm.Open("sqlite3", "./courses.database")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer coursesDB.Close()
	coursesDB.AutoMigrate(&server.Course{})

	if len(os.Args) > 1 {
		seed := os.Args[1]
		if seed == "seed" {
			fmt.Println("Seeding courses from course sqlite file:")
			db.LogMode(true)
			seedDB, err := sql.Open("sqlite3", "./CS188")
			if err != nil {
				fmt.Println(err.Error())
			}

			err = SeedCourses(coursesDB, seedDB)
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
		http.ListenAndServe(":"+PORT, server.Routes(db, coursesDB))
	}
}
