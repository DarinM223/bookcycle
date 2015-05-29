package server

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("helloworld"))
var sessionName string

// InitSessions initializes session store options
func InitSessions(_sessionName string) {
	sessionName = _sessionName
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}
	gob.Register(&User{})
}

// CurrentUser retrieves the current user from the session
func CurrentUser(r *http.Request) (User, error) {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		return User{}, errors.New("You are not logged in")
	}
	if user, ok := sess.Values["user"]; ok {
		if user != nil {
			return *user.(*User), nil
		}
		return User{}, errors.New("You are not logged in")
	}
	return User{}, errors.New("You are not logged in")
}

// SetUserInSession sets a user in the session possibly overwriting existing user
func SetUserInSession(r *http.Request, w http.ResponseWriter, user User) error {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		sess, err = store.New(r, sessionName)
		if err != nil {
			return err
		}
	}
	sess.Values["user"] = user
	err = sess.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser logs a user into a session using a validation function to check passwords, etc
func LoginUser(r *http.Request, w http.ResponseWriter, validateFn func() (User, error)) error {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		sess, err = store.New(r, sessionName)
		if err != nil {
			return err
		}
	}
	if _, ok := sess.Values["user"]; ok {
		return errors.New("User is already logged in")
	}

	user, err := validateFn()
	if err != nil {
		return err
	}

	sess.Values["user"] = user
	err = sess.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// LogoutUser logs a user out of a session
func LogoutUser(r *http.Request, w http.ResponseWriter) error {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	if _, ok := sess.Values["user"]; ok {
		delete(sess.Values, "user")
		err := sess.Save(r, w)
		if err != nil {
			return err
		}
	} else {
		return errors.New("You are not logged in")
	}

	return nil
}