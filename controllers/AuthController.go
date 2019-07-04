// Author: James Mallon <jamesmallondev@gmail.com>
// controllers package - contains the system controllers
package controllers

import (
	. "GoAuthorization/libs/layout"
	. "GoAuthorization/libs/session"
	. "GoAuthorization/models"
	. "GoAuthorization/models/dao"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Struct type authController -
type authController struct {
	LayoutHelper
	pageData map[string]interface{}
}

// AuthController function -
func AuthController() *authController {
	return &authController{}
}

// Signup method - receive a request
func (this *authController) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // if request was post process the form info
		this.signupProcess(w, r)
	}
	this.pageData = map[string]interface{}{"PageTitle": "Index"}
	this.Render(w,
		this.pageData,
		"layout.gohtml", "auth/signup.gohtml")
}

// signupProcess method - process a post request form. data
func (this *authController) signupProcess(w http.ResponseWriter, r *http.Request) {
	user := User{ // create a user object with the form data
		Email:    r.FormValue("email"),
		Password: []byte(r.FormValue("password")),
		Role:     r.FormValue("role")}

	if e := UserDAO().Create(&user); e != nil { // check if email is unique
		http.Error(w, e.Error(), http.StatusForbidden)
		return
	}
	this.Login(w, r) // redirect to login without 302 status, to keep the request state
}

// Login method -
func (this *authController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		this.loginProcess(w, r)
	}
	this.Render(w,
		struct{ PageTitle string }{"Index"},
		"layout.gohtml", "auth/login.gohtml")
}

// loginProcess method - process a post login request
func (this *authController) loginProcess(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	// check user exists and retrieve its password
	user, _ := UserDAO().GetByEmail(email)
	if !(len(user.Email) > 0) {
		http.Error(w, "!Username and/or password do not match", http.StatusForbidden)
		return
	}

	// compare the password
	e := bcrypt.CompareHashAndPassword(user.Password, []byte(pass))
	if e != nil {
		http.Error(w, "Username and/or password do not match", http.StatusForbidden)
		return
	}

	// start session and retrieves the session id
	sid := SessionHelper().Start(w, r)
	// store session
	SessionDAO().Create(&Session{SID: sid, Email: user.Email, LastActivity: time.Now()})

	// redirect
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

// Login method -
func (this *authController) Logout(w http.ResponseWriter, r *http.Request) {
	sid := SessionHelper().Close(w, r)
	if len(sid) > 0 {
		SessionDAO().Remove(sid)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
