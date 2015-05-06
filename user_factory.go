package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type UserFactory interface {
	NewEmptyUser() User                                                          // get an empty user
	NewExistingUser(r *http.Request, paramName string, db gorm.DB) (User, error) // get an existing user from an id route parameter
	NewFormUser(r *http.Request) (User, error)                                   // get a new user from a POST form request
}

type MuxUserFactory struct{}

func NewMuxUserFactory() MuxUserFactory {
	return MuxUserFactory{}
}

func (u MuxUserFactory) NewEmptyUser() User {
	return User{
		Firstname: "",
		Lastname:  "",
		Rating:    0.0,
		Email:     "",
		Phone:     0,
	}
}

func (u MuxUserFactory) NewExistingUser(r *http.Request, paramName string, db gorm.DB) (User, error) {
	user_id, has_user := mux.Vars(r)[paramName]
	if !has_user {
		return User{}, errors.New("User is not defined!")
	}
	var user User
	result := db.First(&user, user_id)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func (u MuxUserFactory) NewFormUser(r *http.Request) (User, error) {
	err := r.ParseForm()
	if err != nil {
		return User{}, err
	}
	first_name := r.PostFormValue("first_name")
	last_name := r.PostFormValue("last_name")
	email := r.PostFormValue("email")
	phone, err := strconv.Atoi(r.PostFormValue("phone"))
	if err != nil {
		return User{}, err
	}
	password := r.PostFormValue("password1")
	password_confirm := r.PostFormValue("password2")
	return NewUser(first_name, last_name, email, phone, password, password_confirm)
}
