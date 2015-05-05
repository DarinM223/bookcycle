package main

import (
	"net/http"
)

type BookFactory interface {
	NewFormBook(r *http.Request, userId int) (Book, error)
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
	version := r.PostFormValue("version")
	class := r.PostFormValue("class")
	professor := r.PostFormValue("professor")
	price := r.PostFormValue("price")
	condition := r.PostFormValue("condition")
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
	}, nil
}
