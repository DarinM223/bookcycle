package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type UserFactory interface {
	NewEmptyUser() User
	NewUser(r *http.Request, paramName string) User
}

type MuxUserFactory struct{}

func NewMuxUserFactory() MuxUserFactory {
	return MuxUserFactory{}
}

func (u MuxUserFactory) NewEmptyUser() User {
	return User{
		Username:  "",
		Firstname: "",
		Lastname:  "",
		Rating:    0.0,
		Email:     "",
		Phone:     0,
	}
}

func (u MuxUserFactory) NewUser(r *http.Request, paramName string) User {
	user_id := mux.Vars(r)[paramName]
	_ = user_id
	return User{
		Username:  "testuser",
		Firstname: "Test",
		Lastname:  "User",
		Rating:    4.5,
		Email:     "testuser@test.com",
		Phone:     123456789,
	}
}
