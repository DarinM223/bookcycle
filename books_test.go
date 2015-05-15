package main

import (
	"net/http"
	"testing"
)

var bookTesting BookTesting

func init() {
	bookTesting = NewBookTesting()
}

func TestCreateBook(t *testing.T) {
	// Test that GET new book route has success when not logged in
	request, err := http.NewRequest("GET", bookTesting.NewBookUrl(), nil)
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

	// Test that POST new book fails when not logged in
	request, err = http.NewRequest("POST", bookTesting.NewBookUrl(), nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 401 {
		t.Errorf("POST 401 expected: %d", res.StatusCode)
		return
	}
	// Test that POST new book succeeds when logged in
	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	err = bookTesting.MakeTestUser(test_user, "password", "password")
	var loginCookie *http.Cookie
	loginCookie, err = bookTesting.LoginUser(test_user.Email, "password")

	test_book := Book{
		Title:     "New book",
		Author:    "Test author",
		Version:   1.0,
		Class:     "Horror",
		Professor: "Smallberg",
		Price:     12.50,
		Condition: 5,
		Details:   "Sample text",
		UserId:    1,
	}
	bookTesting.MakeTestBook(test_book, loginCookie)

	// Test that new book is created
	var books []Book
	bookTesting.DB.Preload("books").Find(&books)

	if len(books) != 1 {
		t.Errorf("Books length should be 1, instead: %d", len(books))
		return
	}

	if books[0].Title != "New book" {
		t.Errorf("\"New book\" expected: %s", books[0].Title)
		return
	}
	if books[0].Author != "Test author" {
		t.Errorf("\"Test author\" expected: %s", books[0].Author)
		return
	}
	if books[0].Version != 1.0 {
		t.Errorf("\"1.0\" expected: %f", books[0].Version)
		return
	}
	if books[0].Class != "Horror" {
		t.Errorf("\"Horror\" expected: %s", books[0].Class)
		return
	}
	if books[0].Professor != "Smallberg" {
		t.Errorf("\"Smallberg\" expected: %s", books[0].Professor)
		return
	}
	if books[0].Price != 12.50 {
		t.Errorf("\"12.50\" expected: %f", books[0].Price)
		return
	}
	if books[0].Condition != 5 {
		t.Errorf("\"5\" expected: %d", books[0].Condition)
		return
	}
	if books[0].Details != "Sample text" {
		t.Errorf("\"Sample Text\" expected: %s", books[0].Details)
		return
	}

	bookTesting.DB.Delete(&books[0])

	// Delete book
	var user User
	bookTesting.DB.Where("email LIKE ?", test_user.Email).First(&user)
	bookTesting.DB.Delete(&user)
}
