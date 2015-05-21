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
	request, err := http.NewRequest("GET", bookTesting.NewBookURL(), nil)
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
	request, err = http.NewRequest("POST", bookTesting.NewBookURL(), nil)
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
	testUser := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	err = bookTesting.MakeTestUser(testUser, "password", "password")
	var loginCookie *http.Cookie
	loginCookie, err = bookTesting.LoginUser(testUser.Email, "password")

	testBook := Book{
		Title:     "New book",
		Author:    "Test author",
		Version:   1.0,
		Class:     "Horror",
		Professor: "Smallberg",
		Price:     12.50,
		Condition: 5,
		Details:   "Sample text",
		UserID:    1,
	}
	bookTesting.MakeTestBook(testBook, loginCookie)

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
	bookTesting.DB.Where("email LIKE ?", testUser.Email).First(&user)
	bookTesting.DB.Delete(&user)
}

func TestDeleteBook(t *testing.T) {
	// test that deleting book without being logged in should fail
	request, err := http.NewRequest("GET", bookTesting.DeleteBookURL(1), nil)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 401 {
		t.Errorf("GET 401 expected: %d", res.StatusCode)
		return
	}

	testUser := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	err = bookTesting.MakeTestUser(testUser, "password", "password")
	var loginCookie *http.Cookie
	loginCookie, err = bookTesting.LoginUser(testUser.Email, "password")

	// test that deleting book while being logged in fails if book does not exist
	request.AddCookie(loginCookie)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 401 {
		t.Errorf("GET 401 expected: %d", res.StatusCode)
		return
	}

	testBook := Book{
		Title:     "New book",
		Author:    "Test author",
		Version:   1.0,
		Class:     "Horror",
		Professor: "Smallberg",
		Price:     12.50,
		Condition: 5,
		Details:   "Sample text",
		UserID:    1,
	}
	bookTesting.MakeTestBook(testBook, loginCookie)

	var myBook Book
	bookTesting.DB.Where("title LIKE ?", testBook.Title).First(&myBook)

	// test that deleting book that you do not own when logged in fails
	newTestUser := User{
		Firstname: "New",
		Lastname:  "User",
		Email:     "newuser@gmail.com",
		Phone:     123456789,
	}
	err = bookTesting.MakeTestUser(newTestUser, "password", "password")
	var newLoginCookie *http.Cookie
	newLoginCookie, err = bookTesting.LoginUser(newTestUser.Email, "password")
	request, err = http.NewRequest("GET", bookTesting.DeleteBookURL(myBook.ID), nil)
	request.AddCookie(newLoginCookie)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 401 {
		t.Errorf("GET 401 expected: %d", res.StatusCode)
		return
	}

	// test that deleting book that exists when logged in succeeds
	request, err = http.NewRequest("GET", bookTesting.DeleteBookURL(myBook.ID), nil)
	request.AddCookie(loginCookie)
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

	// test that there is no more books after deleting
	var books []Book
	bookTesting.DB.Preload("books").Find(&books)
	if len(books) != 0 {
		t.Errorf("Books length 0 expected: %d", len(books))
		return
	}

	// Delete mock created user and book
	var user, newUser User
	bookTesting.DB.Where("email LIKE ?", testUser.Email).First(&user)
	bookTesting.DB.Where("email LIKE ?", newTestUser.Email).First(&newUser)
	bookTesting.DB.Delete(&user)
	bookTesting.DB.Delete(&newUser)
}
