package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CoursesJSONHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	courseID := mux.Vars(r)["id"]

	var course Course
	result := db.First(&course, courseID)
	if result.Error != nil {
		http.Error(w, "Course does not exist", http.StatusUnauthorized)
		return
	}

	courseJSON, err := json.Marshal(course)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(courseJSON)
}
