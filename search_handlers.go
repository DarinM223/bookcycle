package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

// SearchCourse helper function for searching courses
func SearchCourse(department string, courseID string, professor string, db gorm.DB) ([]Course, error) {
	var searchCourses []Course
	if len(department) == 0 {
		return []Course{}, nil
	} else if len(courseID) == 0 {
		result := db.Select("DISTINCT department").Where("department LIKE ?", "%"+department+"%").Limit(10).Find(&searchCourses)
		if result.Error != nil {
			return []Course{}, result.Error
		}
		return searchCourses, nil
	} else if len(professor) == 0 {
		result := db.Select("DISTINCT course_id").Where("department LIKE ? AND course_id LIKE ?", "%"+department+"%", "%"+courseID+"%").
			Limit(10).Find(&searchCourses)
		if result.Error != nil {
			return []Course{}, result.Error
		}
		return searchCourses, nil
	}

	var result *gorm.DB
	professorArray := strings.Split(professor, " ")
	if len(professorArray) < 2 {
		result = db.Where(`department LIKE ? 
							AND course_id LIKE ? 
							AND professor_first_name LIKE ?`,
			"%"+department+"%", "%"+courseID+"%", "%"+professorArray[0]+"%").
			Limit(10).Find(&searchCourses)
	} else {
		result = db.Where(`department LIKE ? 
							AND course_id LIKE ? 
							AND professor_first_name LIKE ? 
							AND professor_last_name LIKE ?`,
			"%"+department+"%", "%"+courseID+"%", "%"+professorArray[0]+"%", "%"+professorArray[1]+"%").
			Limit(10).Find(&searchCourses)
	}
	if result.Error != nil {
		return []Course{}, result.Error
	}
	return searchCourses, nil
}

// SearchBook helper function for searching books
func SearchBook(query string, db gorm.DB) ([]Book, error) {
	var searchBooks []Book
	result := db.Where("title LIKE ?", "%"+query+"%").Limit(10).Find(&searchBooks)
	if result.Error != nil {
		return []Book{}, result.Error
	}
	return searchBooks, nil
}

// SearchResultsJSONHandler Route /search_results.json?query=
func SearchResultsJSONHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	searchBooks, err := SearchBook(query, db)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	searchBooksJSON, err := json.Marshal(searchBooks)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(searchBooksJSON)
}

// SearchResultsHandler Route /search_results?query=
func SearchResultsHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	searchBooks, err := SearchBook(query, db)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t, params, err := GenerateFullTemplate(r, "templates/search_results.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, ManyBookTemplateType{
		UserTemplateType: params,
		Books:            searchBooks,
		Title:            "Search Results",
	})
}

// CourseSearchHandler Route /course_search.json?department=&course_id=&professor=
func CourseSearchHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	department := r.URL.Query().Get("department")
	courseID := r.URL.Query().Get("course_id")
	professor := r.URL.Query().Get("professor")

	searchCourses, err := SearchCourse(department, courseID, professor, db)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	searchCoursesJSON, err := json.Marshal(searchCourses)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(searchCoursesJSON)
}
