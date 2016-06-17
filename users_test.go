package main

import (
	"github.com/DarinM223/bookcycle/server"
	_ "io/ioutil"
	"net/http"
	"testing"
)

var userTesting UserTesting

func init() {
	userTesting = NewUserTesting()
}

func TestCreateUser(t *testing.T) {
	request, err := http.NewRequest("GET", userTesting.NewUserURL(), nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	// Test that GET request returns success
	if res.StatusCode != 200 {
		t.Fatalf("GET Success expected: %d", res.StatusCode)
	}

	testUser := server.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}

	// Test that creating User with no password has error
	if err = userTesting.MakeTestUser(testUser, "", ""); err == nil {
		t.Fatal("Creating User with no password should return error")
	}

	// Test that creating User with different passwords has error
	if err = userTesting.MakeTestUser(testUser, "password1", "password2"); err == nil {
		t.Fatal("Creating User with wrong passwords should return error")
	}

	// Test that creating User properly returns success
	if err = userTesting.MakeTestUser(testUser, "password", "password"); err != nil {
		t.Fatal(err)
	}

	// Test if user is created and proper fields are created
	var users []server.User
	userTesting.DB.Preload("users").Find(&users)
	if len(users) != 1 {
		t.Fatalf("Users length should be 1, instead: %d", len(users))
	}

	if users[0].Firstname != "Test" {
		t.Errorf("\"Test\" expected: %d", users[0].Firstname)
	}
	if users[0].Lastname != "User" {
		t.Errorf("\"User\" expected: %d", users[0].Lastname)
	}
	if users[0].Email != "testuser@gmail.com" {
		t.Errorf("\"testuser@gmail.com\" expected: %d", users[0].Email)
	}
	if users[0].Phone != 123456789 {
		t.Errorf("\"123456789\" expected: %d", users[0].Phone)
	}
	if users[0].Password == "password" {
		t.Errorf("Password not hashed properly: %s", users[0].Password)
	}

	userTesting.DB.Delete(&users[0])
}

func TestEditUser(t *testing.T) {
	testUser := server.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	if err := userTesting.MakeTestUser(testUser, "password", "password"); err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("GET", userTesting.EditUserURL(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test that GET request returns 404 if not logged in
	if res, err := http.DefaultClient.Do(request); err != nil || res.StatusCode != 404 {
		t.Fatal("GET 404 expected")
	}

	loginCookie, err := userTesting.LoginUser(testUser.Email, "password")
	if err != nil {
		t.Fatalf("Error logging in user: %s", err.Error())
	}

	request, err = http.NewRequest("GET", userTesting.EditUserURL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	request.AddCookie(loginCookie) // add login session cookie

	// Test that GET request returns success if logged in
	if res, err := http.DefaultClient.Do(request); err != nil || res.StatusCode != 200 {
		t.Fatal("GET 200 expected")
	}

	editedUser := server.User{
		ID:        testUser.ID,
		Firstname: "T",
		Lastname:  "U",
		Email:     "tu@gmail.com",
		Phone:     987654321,
	}

	// sending edit request with different password and password confirmation should have error
	if err = userTesting.EditTestUser(editedUser, loginCookie, "password1", "password2"); err == nil {
		t.Fatal("Editing user with wrong passwords should have error")
	}

	// sending edit request without password should not change password
	if err = userTesting.EditTestUser(editedUser, loginCookie, "", ""); err != nil {
		t.Fatal(err)
	}

	var user server.User
	userTesting.DB.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("password") {
		t.Fatal("Password should not have changed")
	}

	// sending edit request with a different password should change password
	if err = userTesting.EditTestUser(editedUser, loginCookie, "another_password", "another_password"); err != nil {
		t.Fatal(err)
	}
	userTesting.DB.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("another_password") {
		t.Fatal("Password should have changed")
	}

	userTesting.DB.Delete(&user)
}

func TestViewUser(t *testing.T) {
	testUser := server.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	userRoute := userTesting.ViewUserURL(1)

	// Test that accessing GET for nonexisting user returns 404
	request, err := http.NewRequest("GET", userRoute, nil)
	if res, err := http.DefaultClient.Do(request); err != nil || res.StatusCode != 404 {
		t.Fatal("GET 404 expected")
	}

	// Test that accessing GET for existing user returns success
	if err = userTesting.MakeTestUser(testUser, "password", "password"); err != nil {
		t.Fatal(err)
	}

	// Test that accessing GET for existing user returns 200
	if res, err := http.DefaultClient.Do(request); err != nil || res.StatusCode != 200 {
		t.Fatal("GET success expected")
	}

	// Test that accessing POST for existing user returns 404
	request, err = http.NewRequest("POST", userRoute, nil)
	if res, err := http.DefaultClient.Do(request); err != nil || res.StatusCode != 404 {
		t.Fatal("GET 404 expected")
	}
}
