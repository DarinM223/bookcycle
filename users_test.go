package main

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	db       gorm.DB
	server   *httptest.Server
	usersUrl string
)

func init() {
	// Set up database
	db, _ = gorm.Open("sqlite3", "./sqlite_file_test.db")
	db.LogMode(false)
	db.DropTable(&User{})
	db.DropTable(&Book{})
	db.DropTable(&Book{})
	db.AutoMigrate(&User{}, &Book{}, &Message{})

	// set up test db
	server = httptest.NewServer(Routes(db))

	usersUrl = fmt.Sprintf("%s/users/new", server.URL)
	fmt.Println(usersUrl)
}

func TestCreateUser(t *testing.T) {
	userJson := url.Values{}
	userJson.Set("first_name", "Test")
	userJson.Set("last_name", "User")
	userJson.Set("email", "testuser@gmail.com")
	userJson.Set("phone", "123456789")
	userJson.Set("password1", "password")
	userJson.Set("password2", "password")

	request, err := http.NewRequest("GET", usersUrl, nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	// Test that GET request returns success
	if res.StatusCode != 200 {
		t.Errorf("GET Success expected: %d", res.StatusCode)
		return
	}

	request, err = http.NewRequest("POST", usersUrl, bytes.NewBufferString(userJson.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	// Test that POST request returns success
	if res.StatusCode != 200 {
		t.Errorf("POST Success expected: %d", res.StatusCode)
		return
	}

	// Test if user is created and proper fields are created
	var users []User
	db.Preload("users").Find(&users)
	if len(users) != 1 {
		t.Errorf("Users length should be 1, instead: %d", len(users))
		return
	}

	if users[0].Firstname != "Test" {
		t.Errorf("\"Test\" expected: %d", users[0].Firstname)
		return
	}
	if users[0].Lastname != "User" {
		t.Errorf("\"User\" expected: %d", users[0].Lastname)
		return
	}
	if users[0].Email != "testuser@gmail.com" {
		t.Errorf("\"testuser@gmail.com\" expected: %d", users[0].Email)
		return
	}
	if users[0].Phone != 123456789 {
		t.Errorf("\"123456789\" expected: %d", users[0].Phone)
		return
	}
	if users[0].Password == "password" {
		t.Errorf("Password not hashed properly: %s", users[0].Password)
		return
	}
}
