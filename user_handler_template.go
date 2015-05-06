package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
)

type UserHandler interface {
	Handler(w http.ResponseWriter, r *http.Request, db gorm.DB) // http handler for the specific user route
}

// Hidden interface inside UserHandlerTemplate for doing dynamic method dispatch
type userHandlerTemplateVirtualMethods interface {
	user(r *http.Request, db gorm.DB) (User, error)
	getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
	postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
	isDisabled() bool
}

// Implementation of UserHandler
type UserHandlerTemplate struct {
	userFactory UserFactory
	i           userHandlerTemplateVirtualMethods
}

type UserDetailTemplateType struct {
	DisabledText string
	Disabled     bool
	User         User
}

func (u *UserHandlerTemplate) getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	user, err := u.i.user(r, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	t, err := template.ParseFiles("templates/user_detail.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	disabledText := ""
	if u.i.isDisabled() {
		disabledText = "disabled"
	}
	t.Execute(w, UserDetailTemplateType{
		DisabledText: disabledText,
		Disabled:     u.i.isDisabled(),
		User:         user,
	})
}

func (u UserHandlerTemplate) Handler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		u.i.getRoute(w, r, db)
	} else if r.Method == "POST" {
		u.i.postRoute(w, r, db)
	} else {
		http.NotFound(w, r)
	}
}

// User handler for /users/new
type UserNewTemplate struct {
	UserHandlerTemplate
}

func NewUserNewTemplate() UserNewTemplate {
	b := UserNewTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserNewTemplate) isDisabled() bool { return false }

func (u UserNewTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	return u.userFactory.NewEmptyUser(), nil
}

func (u *UserNewTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	new_user, err := u.userFactory.NewFormUser(r)
	fmt.Println(new_user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	result := db.Create(&new_user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// User handler for route /users/edit
type UserEditTemplate struct {
	UserHandlerTemplate
}

func NewUserEditTemplate() UserEditTemplate {
	b := UserEditTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserEditTemplate) isDisabled() bool { return false }

func (u UserEditTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	user, err := CurrentUser(r)
	if err != nil {
		return User{}, errors.New("You are not logged in")
	}
	return user, nil
}

func (u UserEditTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	// TODO: implement edit existing user
}

// User handler for route /users/{id}
type UserViewTemplate struct {
	UserHandlerTemplate
}

func NewUserViewTemplate() UserViewTemplate {
	b := UserViewTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserViewTemplate) isDisabled() bool { return true }

func (u UserViewTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	user, err := u.userFactory.NewExistingUser(r, "id", db)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u UserViewTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {} // do nothing
