package server

import (
	"net/http"
	"strconv"
	"time"
)

// BookFactory is an interface for creating books from various parameters
type BookFactory interface {
	NewFormBook(r *http.Request, userID int) (Book, error) // Generates new Books from a POST form request
}

// MuxBookFactory is an implementation of BookFactory
type MuxBookFactory struct{}

// NewMuxBookFactory constructs a new MuxBookFactory
func NewMuxBookFactory() MuxBookFactory {
	return MuxBookFactory{}
}

// NewFormBook creates a new book object from a http post form request
func (u MuxBookFactory) NewFormBook(r *http.Request, userID int) (Book, error) {
	if err := r.ParseForm(); err != nil {
		return Book{}, err
	}
	isbn := r.PostFormValue("isbn")
	title := r.PostFormValue("title")
	courseID, err := strconv.Atoi(r.PostFormValue("course_id"))
	if err != nil {
		return Book{}, err
	}
	price, err := strconv.ParseFloat(r.PostFormValue("price"), 64)
	if err != nil {
		return Book{}, err
	}
	condition, err := strconv.Atoi(r.PostFormValue("condition"))
	if err != nil {
		return Book{}, err
	}
	details := r.PostFormValue("details")

	return Book{
		Title:     title,
		ISBN:      isbn,
		CourseID:  courseID,
		Price:     price,
		Condition: condition,
		Details:   details,
		UserID:    userID,
		CreatedAt: time.Now(),
	}, nil
}
