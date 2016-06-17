package server

import (
	"html/template"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/justinas/nosurf"
)

// UserHandler abstracts an http handler with an included database parameter
type UserHandler interface {
	Handler(w http.ResponseWriter, r *http.Request, db gorm.DB) // http handler for the specific user route
}

// Hidden interface inside UserHandlerTemplate for doing dynamic method dispatch
type userHandlerTemplateVirtualMethods interface {
	user(r *http.Request, db gorm.DB) (User, error)
	getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
	postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB)
	isDisabled() bool
}

// UserHandlerTemplate is an implementation of UserHandler
type UserHandlerTemplate struct {
	userFactory UserFactory
	i           userHandlerTemplateVirtualMethods
}

// UserDetailTemplateType abstracts the user detail template
type UserDetailTemplateType struct {
	DisabledText   string
	Disabled       bool
	User           User
	CurrentUser    User
	HasCurrentUser bool
	Token          string
}

func (u *UserHandlerTemplate) getRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	hasCurrentUser := err == nil
	user, err := u.i.user(r, db)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("templates/boilerplate/navbar_boilerplate.html",
		"templates/navbar.html", "templates/user_detail.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	disabledText := ""
	if u.i.isDisabled() {
		disabledText = "disabled"
	}
	t.Execute(w, UserDetailTemplateType{
		DisabledText:   disabledText,
		Disabled:       u.i.isDisabled(),
		User:           user,
		CurrentUser:    currentUser,
		HasCurrentUser: hasCurrentUser,
		Token:          nosurf.Token(r),
	})
}

// Handler abstracts the Http handler and calls either the virtual getRoute or postRoute
// depending on method
func (u UserHandlerTemplate) Handler(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	if r.Method == "GET" {
		u.i.getRoute(w, r, db)
	} else if r.Method == "POST" {
		u.i.postRoute(w, r, db)
	} else {
		http.NotFound(w, r)
	}
}

// UserNewTemplate User handler for /users/new
type UserNewTemplate struct {
	UserHandlerTemplate
}

// NewUserNewTemplate constructs a new UserNewTemplate
func NewUserNewTemplate() UserNewTemplate {
	b := UserNewTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserNewTemplate) isDisabled() bool { return false }

func (u UserNewTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	return u.userFactory.NewEmptyUser(), nil
}

func (u *UserNewTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	newUser, err := u.userFactory.NewFormUser(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if result := db.Create(&newUser); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// UserEditTemplate is a user handler for route /users/edit
type UserEditTemplate struct {
	UserHandlerTemplate
}

// NewUserEditTemplate constructs a new UserEditTemplate
func NewUserEditTemplate() UserEditTemplate {
	b := UserEditTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserEditTemplate) isDisabled() bool { return false }

func (u UserEditTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	return CurrentUser(r)
}

func (u UserEditTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {
	currentUser, err := CurrentUser(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	editedUser, err := u.userFactory.NewFormUser(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if result := db.Model(&currentUser).Updates(editedUser); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	// get edited user from database
	var newUser User
	if result := db.First(&newUser, currentUser.ID); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		return
	}
	if err = SetUserInSession(r, w, newUser); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// UserViewTemplate is a user handler for route /users/{id}
type UserViewTemplate struct {
	UserHandlerTemplate
}

// NewUserViewTemplate constructs a new UserViewTemplate
func NewUserViewTemplate() UserViewTemplate {
	b := UserViewTemplate{UserHandlerTemplate{}}
	b.userFactory = NewMuxUserFactory()
	b.i = &b
	return b
}

func (u UserViewTemplate) isDisabled() bool { return true }

func (u UserViewTemplate) user(r *http.Request, db gorm.DB) (User, error) {
	return u.userFactory.NewExistingUser(r, "id", db)
}

func (u UserViewTemplate) postRoute(w http.ResponseWriter, r *http.Request, db gorm.DB) {} // do nothing
