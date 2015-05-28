package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// UserFactory is an interface to create users from different parameters
type UserFactory interface {
	NewEmptyUser() User                                                          // get an empty user
	NewExistingUser(r *http.Request, paramName string, db gorm.DB) (User, error) // get an existing user from an id route parameter
	NewFormUser(r *http.Request, editing bool) (User, error)                     // get a new user from a POST form request
}

// MuxUserFactory is an implementation of UserFactory
type MuxUserFactory struct{}

// NewMuxUserFactory constructs a new MuxUserFactory
func NewMuxUserFactory() MuxUserFactory {
	return MuxUserFactory{}
}

// NewEmptyUser creates a user object with all empty properties
func (u MuxUserFactory) NewEmptyUser() User {
	return User{
		Firstname: "",
		Lastname:  "",
		Rating:    0.0,
		Email:     "",
		Phone:     0,
	}
}

// NewExistingUser creates a user object from a user id route parameter
func (u MuxUserFactory) NewExistingUser(r *http.Request, paramName string, db gorm.DB) (User, error) {
	userID, hasUser := mux.Vars(r)[paramName]
	if !hasUser {
		return User{}, errors.New("User is not defined!")
	}
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

// NewFormUser creates a user object from a http post form request
func (u MuxUserFactory) NewFormUser(r *http.Request, editing bool) (User, error) {
	err := r.ParseForm()
	if err != nil {
		return User{}, err
	}
	firstName := r.PostFormValue("first_name")
	lastName := r.PostFormValue("last_name")
	email := r.PostFormValue("email")
	phone, err := strconv.Atoi(r.PostFormValue("phone"))
	if err != nil {
		return User{}, err
	}
	password := r.PostFormValue("password1")
	passwordConfirm := r.PostFormValue("password2")
	return NewUser(firstName, lastName, email, phone, password, passwordConfirm, editing)
}
