package main

import (
	"html/template"
	"net/http"
)

type Routes interface {
	User(r *http.Request) User
	GetRoute(w http.ResponseWriter, r *http.Request)
	PostRoute(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	Handler(w http.ResponseWriter, r *http.Request)
	Error() error
}

type UserHandlerTemplate struct {
	userFactory UserFactory
	err         error
	i           Routes
}

func (u UserHandlerTemplate) Error() error {
	return u.err
}

func (u UserHandlerTemplate) User(r *http.Request) User {
	return u.userFactory.NewEmptyUser()
}

func (u *UserHandlerTemplate) GetRoute(w http.ResponseWriter, r *http.Request) {
	user := u.i.User(r)
	t, err := template.ParseFiles("templates/user_detail.html")
	if err != nil {
		u.err = err
		return
	}
	t.Execute(w, user)
}

func (u UserHandlerTemplate) PostRoute(w http.ResponseWriter, r *http.Request) {}

func (u UserHandlerTemplate) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		u.i.GetRoute(w, r)
	} else if r.Method == "POST" {
		u.i.PostRoute(w, r)
	} else {
		http.NotFound(w, r)
	}
}

type UserNewTemplate struct {
	UserHandlerTemplate
}

func NewUserNewTemplate() UserNewTemplate {
	b := UserNewTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.err = nil
	b.i = &b
	return b
}

func (u UserNewTemplate) User(r *http.Request) User {
	return u.userFactory.NewEmptyUser()
}

func (u UserNewTemplate) PostRoute(w http.ResponseWriter, r *http.Request) {
	// TODO: implement add new user
}

type UserEditTemplate struct {
	UserHandlerTemplate
}

func NewUserEditTemplate() UserEditTemplate {
	b := UserEditTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.err = nil
	b.i = &b
	return b
}

func (u UserEditTemplate) User(r *http.Request) User {
	return u.userFactory.NewUser(r, "id")
}

func (u UserEditTemplate) PostRoute(w http.ResponseWriter, r *http.Request) {
	// TODO: implement edit existing user
}
