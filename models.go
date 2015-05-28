package main

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User has the fields of a user
type User struct {
	ID        int       `sql:"AUTO_INCREMENT" json:"id"`
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

// NewUser constructs a new User
func NewUser(firstname string, lastname string, email string,
	phone int, password string, passwordConfirm string, editing bool) (User, error) {

	if password != passwordConfirm {
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
		}
		return User{}, errors.New("Password cannot be empty")
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	return User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
		Password:  string(encryptedPassword),
		CreatedAt: time.Now(),
	}, nil
}

// Validate validates if the password matches the user's hashed password
func (u User) Validate(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// Book represents a book
type Book struct {
	ID        int       `sql:"AUTO_INCREMENT" json:"id"`
	Title     string    `sql:"not null" json:"title"`
	ISBN      string    `sql:"not null" json:"isbn"`
	Price     float64   `sql:"not null" json:"price"`
	Condition int       `sql:"not null" json:"condition"`
	Details   string    `json:"details"`
	UserID    int       `sql:"index" json:"user_id"`
	CourseID  int       `sql:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// Course represents a UCLA class
type Course struct {
	ID         int       `sql:"AUTO_INCREMENT" json:"id"`
	Department string    `json:"department"`
	CourseID   string    `json:"course_id"`
	Professor  string    `json:"professor"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Message represents a message
type Message struct {
	SenderID   int       `json:"senderId"`
	ReceiverID int       `sql:"index" json:"receiverId"`
	Message    string    `json:"message"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"created_at"`
}
