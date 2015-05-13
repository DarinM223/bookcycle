package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

var (
	db           gorm.DB
	server       *httptest.Server
	newUserUrl   string
	editUserUrl  string
	loginUserUrl string
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

	newUserUrl = fmt.Sprintf("%s/users/new", server.URL)
	editUserUrl = fmt.Sprintf("%s/users/edit", server.URL)
	loginUserUrl = fmt.Sprintf("%s/login", server.URL)
}

func makeTestUser(u User, password string, password_confirm string) error {
	userJson := url.Values{}
	userJson.Set("first_name", u.Firstname)
	userJson.Set("last_name", u.Lastname)
	userJson.Set("email", u.Email)
	userJson.Set("phone", strconv.Itoa(u.Phone))
	userJson.Set("password1", password)
	userJson.Set("password2", password_confirm)

	request, err := http.NewRequest("POST", newUserUrl, bytes.NewBufferString(userJson.Encode()))
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
		//body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(string(body))
		return errors.New("POST Success should be 200")
	}

	return nil
}

func editTestUser(u User, c *http.Cookie, password string, password_confirm string) error {
	userJson := url.Values{}
	userJson.Set("first_name", u.Firstname)
	userJson.Set("last_name", u.Lastname)
	userJson.Set("email", u.Email)
	userJson.Set("phone", strconv.Itoa(u.Phone))
	userJson.Set("password1", password)
	userJson.Set("password2", password_confirm)

	request, err := http.NewRequest("POST", editUserUrl, bytes.NewBufferString(userJson.Encode()))
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
		//body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(string(body))
		return errors.New("POST Success should be 200")
	}

	return nil
}

func loginUser(email string, password string) (*http.Cookie, error) {
	loginJson := url.Values{}
	loginJson.Set("email", email)
	loginJson.Set("password", password)

	request, err := http.NewRequest("POST", loginUserUrl, bytes.NewBufferString(loginJson.Encode()))
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

func TestCreateUser(t *testing.T) {
	request, err := http.NewRequest("GET", newUserUrl, nil)
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

	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}

	// Test that creating User with no password has error
	err = makeTestUser(test_user, "", "")
	if err == nil {
		t.Error("Creating User with no password should return error")
		return
	}

	// Test that creating User with different passwords has error
	err = makeTestUser(test_user, "password1", "password2")
	if err == nil {
		t.Error("Creating User with wrong passwords should return error")
		return
	}

	// Test that creating User properly returns success
	err = makeTestUser(test_user, "password", "password")
	if err != nil {
		t.Error(err)
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

	db.Delete(&users[0])
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func TestEditUser(t *testing.T) {
	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	err := makeTestUser(test_user, "password", "password")
	if err != nil {
		t.Error(err)
		return
	}

	request, err := http.NewRequest("GET", editUserUrl, nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	// Test that GET request returns 404 if not logged in
	if res.StatusCode != 404 {
		t.Errorf("GET 404 expected: %d", res.StatusCode)
		return
	}

	loginCookie, err := loginUser(test_user.Email, "password")
	if err != nil {
		t.Error(err)
		t.Error("Error logging in user")
		return
	}

	request, err = http.NewRequest("GET", editUserUrl, nil)
	request.AddCookie(loginCookie) // add login session cookie
	if err != nil {
		t.Error(err)
		return
	}

	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	// Test that GET request returns success if logged in
	if res.StatusCode != 200 {
		t.Errorf("GET success expected: %d", res.StatusCode)
		return
	}

	edited_user := User{
		Id:        test_user.Id,
		Firstname: "T",
		Lastname:  "U",
		Email:     "tu@gmail.com",
		Phone:     987654321,
	}

	// sending edit request with different password and password confirmation should have error
	err = editTestUser(edited_user, loginCookie, "password1", "password2")
	if err == nil {
		t.Error("Editing user with wrong passwords should have error")
		return
	}

	// sending edit request without password should not change password
	err = editTestUser(edited_user, loginCookie, "", "")
	if err != nil {
		t.Error(err)
		return
	}

	var user User
	db.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("password") {
		t.Error("Password should not have changed")
		return
	}

	// sending edit request with a different password should change password
	err = editTestUser(edited_user, loginCookie, "another_password", "another_password")
	db.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("another_password") {
		t.Error("Password should have changed")
		return
	}

	db.Delete(&user)
}

func TestViewUser(t *testing.T) {
	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	user_route := fmt.Sprintf("%s/users/1", server.URL)

	// Test that accessing GET for nonexisting user returns 404
	request, err := http.NewRequest("GET", user_route, nil)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != 404 {
		t.Errorf("GET 404 expected: %d", res.StatusCode)
		return
	}

	// Test that accessing GET for existing user returns success
	err = makeTestUser(test_user, "password", "password")
	if err != nil {
		t.Error(err)
		return
	}
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != 200 {
		t.Errorf("GET success expected: %d", res.StatusCode)
		return
	}

	// Test that accessing POST for existing user returns 404
	request, err = http.NewRequest("POST", user_route, nil)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	// Test that POST request returns 404
	if res.StatusCode != 404 {
		t.Errorf("GET 404 expected: %d", res.StatusCode)
		return
	}
}
