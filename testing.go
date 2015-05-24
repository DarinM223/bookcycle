package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/jinzhu/gorm"
)

// SetUpTesting starts a test http server and sets up a test database
func SetUpTesting() (*httptest.Server, gorm.DB) {
	// Set up database
	db, _ := gorm.Open("sqlite3", "./sqlite_file_test.db")
	db.LogMode(false)
	db.DropTable(&User{})
	db.DropTable(&Book{})
	db.DropTable(&Book{})
	db.AutoMigrate(&User{}, &Book{}, &Message{})

	// set up test db
	server := httptest.NewServer(Routes(db))

	return server, db
}

// UserTesting is a struct for User testing utilities
type UserTesting struct {
	DB     gorm.DB
	Server *httptest.Server
}

// NewUserTesting constructs a new UserTesting
func NewUserTesting() UserTesting {
	server, db := SetUpTesting()
	return UserTesting{
		DB:     db,
		Server: server,
	}
}

// NewUserURL returns the new user url
func (n UserTesting) NewUserURL() string {
	return fmt.Sprintf("%s/users/new", n.Server.URL)
}

// EditUserURL returns the edit user url
func (n UserTesting) EditUserURL() string {
	return fmt.Sprintf("%s/users/edit", n.Server.URL)
}

// LoginUserURL returns the login user url
func (n UserTesting) LoginUserURL() string {
	return fmt.Sprintf("%s/login", n.Server.URL)
}

// ViewUserURL returns the view user url
func (n UserTesting) ViewUserURL(id int) string {
	return fmt.Sprintf("%s/users/%d", n.Server.URL, id)
}

// MakeTestUser makes a new test user
func (n UserTesting) MakeTestUser(u User, password string, passwordConfirm string) error {
	userJSON := url.Values{}
	userJSON.Set("first_name", u.Firstname)
	userJSON.Set("last_name", u.Lastname)
	userJSON.Set("email", u.Email)
	userJSON.Set("phone", strconv.Itoa(u.Phone))
	userJSON.Set("password1", password)
	userJSON.Set("password2", passwordConfirm)

	request, err := http.NewRequest("POST", n.NewUserURL(), bytes.NewBufferString(userJSON.Encode()))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	// Test that POST request returns success
	if res.StatusCode != 200 {
		return errors.New("POST Success should be 200")
	}

	return nil
}

// EditTestUser edits an existing user
func (n UserTesting) EditTestUser(u User, c *http.Cookie, password string, passwordConfirm string) error {
	userJSON := url.Values{}
	userJSON.Set("first_name", u.Firstname)
	userJSON.Set("last_name", u.Lastname)
	userJSON.Set("email", u.Email)
	userJSON.Set("phone", strconv.Itoa(u.Phone))
	userJSON.Set("password1", password)
	userJSON.Set("password2", passwordConfirm)

	request, err := http.NewRequest("POST", n.EditUserURL(), bytes.NewBufferString(userJSON.Encode()))
	request.AddCookie(c)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	// Test that POST request returns success
	if res.StatusCode != 200 {
		return errors.New("POST Success should be 200")
	}

	return nil
}

// LoginUser logs in a user
func (n UserTesting) LoginUser(email string, password string) (*http.Cookie, error) {
	loginJSON := url.Values{}
	loginJSON.Set("email", email)
	loginJSON.Set("password", password)

	request, err := http.NewRequest("POST", n.LoginUserURL(), bytes.NewBufferString(loginJSON.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	transport := http.Transport{}
	res, err := transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "bookcycle" {
			return cookie, nil
		}
	}

	return nil, errors.New("Cookie not set")
}

// BookTesting is a struct for book testing utilties
type BookTesting struct {
	UserTesting
}

// NewBookTesting constructs a new BookTesting
func NewBookTesting() BookTesting {
	server, db := SetUpTesting()
	return BookTesting{UserTesting{DB: db, Server: server}}
}

// ShowBooksURL returns the show books url
func (b BookTesting) ShowBooksURL() string {
	return fmt.Sprintf("%s/books", b.Server.URL)
}

// NewBookURL returns the new book url
func (b BookTesting) NewBookURL() string {
	return fmt.Sprintf("%s/books/new", b.Server.URL)
}

// DeleteBookURL returns the delete book url
func (b BookTesting) DeleteBookURL(id int) string {
	return fmt.Sprintf("%s/books/%d/delete", b.Server.URL, id)
}

// ShowBookURL returns the show book url
func (b BookTesting) ShowBookURL(id int) string {
	return fmt.Sprintf("%s/books/%d", b.Server.URL, id)
}

// MakeTestBook makes a new test book
func (b BookTesting) MakeTestBook(book Book, loginCookie *http.Cookie) error {
	bookJSON := url.Values{}
	bookJSON.Set("title", book.Title)
	bookJSON.Set("isbn", book.ISBN)
	bookJSON.Set("course_id", strconv.Itoa(book.CourseID))
	bookJSON.Set("price", fmt.Sprintf("%f", book.Price))
	bookJSON.Set("condition", strconv.Itoa(book.Condition))
	bookJSON.Set("details", book.Details)

	request, err := http.NewRequest("POST", b.NewBookURL(), bytes.NewBufferString(bookJSON.Encode()))
	request.AddCookie(loginCookie)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	// Test that POST request returns success
	if res.StatusCode != 200 {
		return errors.New("POST Success should be 200")
	}

	return nil
}
