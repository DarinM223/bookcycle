package main

import (
	"net/http"
	"strconv"
	"time"
)

type BookFactory interface {
	NewFormBook(r *http.Request, userId int) (Book, error) // Generates new Books from a POST form request
}

type MuxBookFactory struct{}

func NewMuxBookFactory() MuxBookFactory {
	return MuxBookFactory{}
}

func (u MuxBookFactory) NewFormBook(r *http.Request, userId int) (Book, error) {
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
		UserId:    userId,
		CreatedAt: time.Now(),
	}, nil
}
