package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	Id        int     `sql:"AUTO_INCREMENT"`
	Firstname string  `sql:"not null"`
	Lastname  string  `sql:"not null"`
	Rating    float64 `sql:"not null; default:0"`
	Email     string  `sql:"not null; unique"`
	Phone     int
	Password  string `sql:"not null"`
	Books     []Book
}

func NewUser(firstname string, lastname string, email string,
	phone int, password string, password_confirm string) (User, error) {

	if password != password_confirm {
		return User{}, errors.New("Passwords do not match")
	}
	if len(strings.Trim(password, " ")) == 0 {
		return User{}, errors.New("Password cannot be empty")
	}
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	return User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
		Password:  string(encrypted_password),
	}, nil
}

func (u User) Validate(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

type Book struct {
	Id        int    `sql:"AUTO_INCREMENT"`
	Title     string `sql:"not null"`
	Author    string `sql:"not null"`
	Class     string `sql:"not null"`
	Professor string `sql:"not null"`
	Version   string `sql:"not null"`
	Price     string `sql:"not null"`
	Condition string `sql:"not null"`
	Details   string
	UserId    int `sql:"index"`
}
