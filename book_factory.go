package main

import (
	"net/http"
	"strconv"
	"time"
)

// BookFactory is an interface for createing books from various parameters
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
	err := r.ParseForm()
	if err != nil {
		return Book{}, err
	}
	title := r.PostFormValue("title")
	author := r.PostFormValue("author")
	version, err := strconv.ParseFloat(r.PostFormValue("version"), 64)
	if err != nil {
		return Book{}, err
	}
	class := r.PostFormValue("class")
	professor := r.PostFormValue("professor")
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
		Author:    author,
		Version:   version,
		Class:     class,
		Professor: professor,
		Price:     price,
		Condition: condition,
		Details:   details,
		UserID:    userID,
		CreatedAt: time.Now(),
	}, nil
}
