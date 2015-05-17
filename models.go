package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	Id        int       `sql:"AUTO_INCREMENT" json:"id"`
	Firstname string    `sql:"not null" json:"first_name"`
	Lastname  string    `sql:"not null" json:"last_name"`
	Rating    float64   `sql:"not null; default:0" json:"rating"`
	Email     string    `sql:"not null; unique" json:"email"`
	Phone     int       `json:"phone"`
	Password  string    `sql:"not null" json:"-"`
	Messages  []Message `json:"-"`
	Books     []Book    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(firstname string, lastname string, email string,
	phone int, password string, password_confirm string, editing bool) (User, error) {

	if password != password_confirm {
		return User{}, errors.New("Passwords do not match")
	}
	if len(strings.Trim(password, " ")) == 0 {
		if editing { // ignore empty passwords if editing
			return User{
				Firstname: firstname,
				Lastname:  lastname,
				Email:     email,
				Phone:     phone,
				UpdatedAt: time.Now(),
			}, nil
		} else {
			return User{}, errors.New("Password cannot be empty")
		}
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
		CreatedAt: time.Now(),
	}, nil
}

func (u User) Validate(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

type Book struct {
	Id        int       `sql:"AUTO_INCREMENT" json:"id"`
	Title     string    `sql:"not null" json:"title"`
	Author    string    `sql:"not null" json:"author"`
	Class     string    `sql:"not null" json:"class"`
	Professor string    `sql:"not null" json:"professor"`
	Version   float64   `sql:"not null" json:"version"`
	Price     float64   `sql:"not null" json:"price"`
	Condition int       `sql:"not null" json:"condition"`
	Details   string    `json:"details"`
	UserId    int       `sql:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	SenderId   int       `json:"senderId"`
	ReceiverId int       `sql:"index" json:"receiverId"`
	Message    string    `json:"message"`
	Read       bool      `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}
