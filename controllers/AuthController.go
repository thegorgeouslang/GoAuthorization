// Author: James Mallon <jamesmallondev@gmail.com>
// controllers package - contains the system controllers
package controllers

import (
	. "GoAuthorization/controllers/helpers"
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
	LayoutManager
	flash FlashMessenger
}

// to fill with the flash message values
type fm map[string]string

// AuthController function -
func AuthController() *authController {
	return &authController{}
}

// Signup method - receive a request
func (this *authController) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // if request was post process the form info
		// filtering form inputs
		if formErrs := FormHelper.Filter(r); len(formErrs) > 0 {
			this.flash.Set(&w, fm{"message": FormHelper.ErrString(formErrs), "type": "danger"})
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		user := User{ // create a user object with the form data
			Email:    r.FormValue("email"),
			Password: []byte(r.FormValue("password")),
			Role:     r.FormValue("role")}

		if e := UserDAO().Create(&user); e != nil { // check if email is unique
			this.flash.Set(&w, fm{"message": e.Error(), "type": "danger"})
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		this.Login(w, r) // redirect to login without 302 status, to keep the request state
	}
	PageData["PageTitle"] = "Signup"
	this.Render(w, r,
		PageData,
		"layout.gohtml", "auth/signup.gohtml")
}

// Login method -
func (this *authController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		pass := r.FormValue("password")

		// check user exists and retrieve its password
		user, _ := UserDAO().GetByEmail(email)
		if !(len(user.Email) > 0) {
			this.flash.Set(&w, fm{"message": "Username and/or password do not match", "type": "danger"})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// compare the password
		e := bcrypt.CompareHashAndPassword(user.Password, []byte(pass))
		if e != nil {
			this.flash.Set(&w, fm{"message": "Username and/or password do not match", "type": "danger"})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// start session and retrieves the session id
		sid := SessionManager().Start(w, r)
		// store session
		SessionDAO().Create(&Session{SID: sid, Email: user.Email, LastActivity: time.Now()})

		// redirect
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
		//	this.loginProcess(w, r)
	}
	PageData["PageTitle"] = "Login"
	this.Render(w, r,
		PageData,
		"layout.gohtml", "auth/login.gohtml")
}

// Login method -
func (this *authController) Logout(w http.ResponseWriter, r *http.Request) {
	sid := SessionManager().Close(w, r)
	if len(sid) > 0 {
		SessionDAO().Remove(sid)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
