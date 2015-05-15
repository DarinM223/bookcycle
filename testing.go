package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
)

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

// Struct for User testing utilities
type UserTesting struct {
	DB     gorm.DB
	Server *httptest.Server
}

func NewUserTesting() UserTesting {
	server, db := SetUpTesting()
	return UserTesting{
		DB:     db,
		Server: server,
	}
}

func (n UserTesting) NewUserUrl() string {
	return fmt.Sprintf("%s/users/new", n.Server.URL)
}

func (n UserTesting) EditUserUrl() string {
	return fmt.Sprintf("%s/users/edit", n.Server.URL)
}

func (n UserTesting) LoginUserUrl() string {
	return fmt.Sprintf("%s/login", n.Server.URL)
}

func (n UserTesting) ViewUserUrl(id int) string {
	return fmt.Sprintf("%s/users/%d", n.Server.URL, id)
}

func (n UserTesting) MakeTestUser(u User, password string, password_confirm string) error {
	userJson := url.Values{}
	userJson.Set("first_name", u.Firstname)
	userJson.Set("last_name", u.Lastname)
	userJson.Set("email", u.Email)
	userJson.Set("phone", strconv.Itoa(u.Phone))
	userJson.Set("password1", password)
	userJson.Set("password2", password_confirm)

	request, err := http.NewRequest("POST", n.NewUserUrl(), bytes.NewBufferString(userJson.Encode()))
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

func (n UserTesting) EditTestUser(u User, c *http.Cookie, password string, password_confirm string) error {
	userJson := url.Values{}
	userJson.Set("first_name", u.Firstname)
	userJson.Set("last_name", u.Lastname)
	userJson.Set("email", u.Email)
	userJson.Set("phone", strconv.Itoa(u.Phone))
	userJson.Set("password1", password)
	userJson.Set("password2", password_confirm)

	request, err := http.NewRequest("POST", n.EditUserUrl(), bytes.NewBufferString(userJson.Encode()))
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

func (n UserTesting) LoginUser(email string, password string) (*http.Cookie, error) {
	loginJson := url.Values{}
	loginJson.Set("email", email)
	loginJson.Set("password", password)

	request, err := http.NewRequest("POST", n.LoginUserUrl(), bytes.NewBufferString(loginJson.Encode()))
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

// Struct for book testing utilties
type BookTesting struct {
	UserTesting
}

func NewBookTesting() BookTesting {
	server, db := SetUpTesting()
	return BookTesting{UserTesting{DB: db, Server: server}}
}

func (b BookTesting) ShowBooksUrl() string {
	return fmt.Sprintf("%s/books", b.Server.URL)
}

func (b BookTesting) NewBookUrl() string {
	return fmt.Sprintf("%s/books/new", b.Server.URL)
}

func (b BookTesting) DeleteBookUrl(id int) string {
	return fmt.Sprintf("%s/books/%d/delete", b.Server.URL, id)
}

func (b BookTesting) ShowBookUrl(id int) string {
	return fmt.Sprintf("%s/books/%d", b.Server.URL, id)
}

func (b BookTesting) MakeTestBook(book Book, loginCookie *http.Cookie) error {
	bookJson := url.Values{}
	bookJson.Set("title", book.Title)
	bookJson.Set("author", book.Author)
	bookJson.Set("version", fmt.Sprintf("%f", book.Version))
	bookJson.Set("class", book.Class)
	bookJson.Set("professor", book.Professor)
	bookJson.Set("price", fmt.Sprintf("%f", book.Price))
	bookJson.Set("condition", strconv.Itoa(book.Condition))
	bookJson.Set("details", book.Details)

	request, err := http.NewRequest("POST", b.NewBookUrl(), bytes.NewBufferString(bookJson.Encode()))
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
