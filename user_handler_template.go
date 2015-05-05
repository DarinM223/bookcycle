package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
)

type Routes interface {
	user(r *http.Request) User
	getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
	postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
}

type UserHandler interface {
	Handler(w http.ResponseWriter, r *http.Request, db gorm.DB)
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

func (u UserHandlerTemplate) user(r *http.Request) User {
	return u.userFactory.NewEmptyUser()
}

func (u *UserHandlerTemplate) getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	user := u.i.user(r)
	t, err := template.ParseFiles("templates/user_detail.html")
	if err != nil {
		u.err = err
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	t.Execute(w, user)
}

func (u UserHandlerTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {}

func (u UserHandlerTemplate) Handler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		u.i.getRoute(w, r, db)
	} else if r.Method == "POST" {
		u.i.postRoute(w, r, db)
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

func (u UserNewTemplate) user(r *http.Request) User {
	return u.userFactory.NewEmptyUser()
}

func (u *UserNewTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	new_user, err := u.userFactory.NewFormUser(r)
	fmt.Println(new_user)
	if err != nil {
		u.err = err
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	result := db.Create(&new_user)
	if result.Error != nil {
		u.err = result.Error
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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

func (u UserEditTemplate) user(r *http.Request) User {
	user, err := u.userFactory.NewUser(r, "id")
	if err != nil {
		user = u.userFactory.NewEmptyUser()
	}
	return user
}

func (u UserEditTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	// TODO: implement edit existing user
}
