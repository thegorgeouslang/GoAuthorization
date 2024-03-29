// Author: James Mallon <jamesmallondev@gmail.com>
// session package -
package session

import (
	conf "GoAuthorization/configs"
	log "GoAuthorization/libs/logger"
	. "GoAuthorization/models"
	. "GoAuthorization/models/dao"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

// Struct type session -
type sessionManager struct {
	sessExp int
}

// init function - data and process initialization
func SessionManager() *sessionManager {
	se, _ := strconv.Atoi(conf.Env["session_exp"])
	return &sessionManager{se}
}

// GetSession method -
func (this *sessionManager) Start(w http.ResponseWriter, r *http.Request) string {
	c, e := r.Cookie("session") // create the cookie
	if e != nil {
		log.Write("notice", "Creating the sesion cookie "+e.Error(), log.Trace())
		sID := uuid.NewV4() // create the universal unique id
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = this.sessExp
	http.SetCookie(w, c)
	return c.Value
}

// Get method -
func (this *sessionManager) GetSession(w http.ResponseWriter, r *http.Request) (e error) {
	c, e := r.Cookie("session") // create the cookie
	if e != nil {
		return
	}
	session, _ := SessionDAO().GetSession(c.Value)
	if session.Email != "" {
		SessionDAO().Renew(c.Value)
	}
	// refresh session
	c.MaxAge = this.sessExp
	http.SetCookie(w, c)
	return
}

//
func (this *sessionManager) User(w http.ResponseWriter, r *http.Request) (user User) {
	c, e := r.Cookie("session") // create the cookie
	if e == nil {
		c.MaxAge = this.sessExp
		http.SetCookie(w, c)

		// if the user exists already, get user
		session, _ := SessionDAO().GetSession(c.Value) // retrieve the session
		if len(session.Email) > 0 {                    // check for the user email
			SessionDAO().Renew(c.Value)                   // update LastActivity
			user, _ = UserDAO().GetByEmail(session.Email) // retrieve user
			return
		}
	}
	return
}

// CloseSession method -
func (this *sessionManager) Close(w http.ResponseWriter, r *http.Request) (sid string) {
	c, err := r.Cookie("session") // create the cookie
	if err == nil {
		sid = c.Value
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	return
}
