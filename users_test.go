package main

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	server   *httptest.Server
	usersUrl string
)

func init() {
	// Set up database
	db, err := gorm.Open("sqlite3", "./sqlite_file_test.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		t.Errorf("POST Success expected: %d", res.StatusCode)
		return
	}
}
