package main

import (
	_ "io/ioutil"
	"net/http"
	"testing"
)

var userTesting UserTesting

func init() {
	userTesting = NewUserTesting()
}

func TestCreateUser(t *testing.T) {
	request, err := http.NewRequest("GET", userTesting.NewUserUrl(), nil)
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
	err = userTesting.MakeTestUser(test_user, "", "")
	if err == nil {
		t.Error("Creating User with no password should return error")
		return
	}

	// Test that creating User with different passwords has error
	err = userTesting.MakeTestUser(test_user, "password1", "password2")
	if err == nil {
		t.Error("Creating User with wrong passwords should return error")
		return
	}

	// Test that creating User properly returns success
	err = userTesting.MakeTestUser(test_user, "password", "password")
	if err != nil {
		t.Error(err)
		return
	}

	// Test if user is created and proper fields are created
	var users []User
	userTesting.DB.Preload("users").Find(&users)
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

	userTesting.DB.Delete(&users[0])
}

func TestEditUser(t *testing.T) {
	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	err := userTesting.MakeTestUser(test_user, "password", "password")
	if err != nil {
		t.Error(err)
		return
	}

	request, err := http.NewRequest("GET", userTesting.EditUserUrl(), nil)
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

	loginCookie, err := userTesting.LoginUser(test_user.Email, "password")
	if err != nil {
		t.Error(err)
		t.Error("Error logging in user")
		return
	}

	request, err = http.NewRequest("GET", userTesting.EditUserUrl(), nil)
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
	err = userTesting.EditTestUser(edited_user, loginCookie, "password1", "password2")
	if err == nil {
		t.Error("Editing user with wrong passwords should have error")
		return
	}

	// sending edit request without password should not change password
	err = userTesting.EditTestUser(edited_user, loginCookie, "", "")
	if err != nil {
		t.Error(err)
		return
	}

	var user User
	userTesting.DB.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("password") {
		t.Error("Password should not have changed")
		return
	}

	// sending edit request with a different password should change password
	err = userTesting.EditTestUser(edited_user, loginCookie, "another_password", "another_password")
	userTesting.DB.Where("email LIKE ?", "tu@gmail.com").First(&user)
	if !user.Validate("another_password") {
		t.Error("Password should have changed")
		return
	}

	userTesting.DB.Delete(&user)
}

func TestViewUser(t *testing.T) {
	test_user := User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@gmail.com",
		Phone:     123456789,
	}
	user_route := userTesting.ViewUserUrl(1)

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
	err = userTesting.MakeTestUser(test_user, "password", "password")
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
