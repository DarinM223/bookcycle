package main

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("helloworld"))
var sessionName string

func InitSessions(_sessionName string) {
	sessionName = _sessionName
	store.Options = &sessions.Options{
		//Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}
	gob.Register(&User{})
}

func CurrentUser(r *http.Request, w http.ResponseWriter) (User, error) {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		return User{}, errors.New("You are not logged in")
	}
	if user, ok := sess.Values["user"]; ok {
		if user != nil {
			return *user.(*User), nil
		} else {
			return User{}, errors.New("You are not logged in")
		}
	} else {
		return User{}, errors.New("You are not logged in")
	}
}

func LoginUser(r *http.Request, w http.ResponseWriter, validateFn func() (User, error)) error {
	sess, err := store.Get(r, "bookcycle")
	if err != nil {
		sess, err = store.New(r, "bookcycle")
		if err != nil {
			return err
		}
	}
	if _, ok := sess.Values["user"]; ok {
		return errors.New("User is already logged in")
	} else {
		user, err := validateFn()
		if err != nil {
			return err
		}

		sess.Values["user"] = user
		err = sess.Save(r, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func LogoutUser(r *http.Request, w http.ResponseWriter) error {
	sess, err := store.Get(r, "bookcycle")
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
