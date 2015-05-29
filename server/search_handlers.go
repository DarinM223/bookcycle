package server

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"net/http"
)

// SearchCourse is a helper function that takes in a search type, department, course id, and professor and returns all courses that match
// search type can be:
// department (when searching for department)
// course (when searching for course)
// professor (when searching for professor)
func SearchCourse(searchType string, department string, courseID string, professor string, db gorm.DB) ([]Course, error) {
	var result *gorm.DB
	var searchCourses []Course

	switch searchType {
	case "department":
		result = db.Select("DISTINCT department").Where("department LIKE ?", "%"+department+"%").Limit(10).Find(&searchCourses)
	case "course":
		result = db.Select("DISTINCT course_id").Where("department LIKE ? AND course_id LIKE ?", department, "%"+courseID+"%").
			Limit(10).Find(&searchCourses)
	case "professor":
		result = db.Where(`department LIKE ? 
						   AND course_id LIKE ? 
						   AND professor LIKE ?`,
			department, courseID, "%"+professor+"%").
			Limit(10).Find(&searchCourses)
	default:
		return []Course{}, errors.New("No search type")
	}

	if result.Error != nil {
		return []Course{}, result.Error
	}
	return searchCourses, nil
}

// SearchResultsJSONHandler is a route for /search_results.json?query= that returns an array of Books that match the search query in JSON format
func SearchResultsJSONHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	var searchBooks []Book
	result := db.Select("DISTINCT title").Where("title LIKE ?", "%"+query+"%").Limit(10).Find(&searchBooks)
	if result.Error != nil {
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

// SearchResultsHandler is a route for /search_results?query= that displays a search page with Books that match the search query
func SearchResultsHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.NotFound(w, r)
		return
	}

	var searchBooks []Book
	result := db.Where("title LIKE ?", "%"+query+"%").Limit(10).Find(&searchBooks)
	if result.Error != nil {
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

// CourseSearchHandler is a route for /course_search.json?department=&course_id=&professor= that returns an array of Courses that match the queries
// GET parameters:
// type string (type of search (department, course, professor))
// department string
// course_id string
// professor string
func CourseSearchHandler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	typeSearch := r.URL.Query().Get("type")
	department := r.URL.Query().Get("department")
	courseID := r.URL.Query().Get("course_id")
	professor := r.URL.Query().Get("professor")

	searchCourses, err := SearchCourse(typeSearch, department, courseID, professor, db)
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
