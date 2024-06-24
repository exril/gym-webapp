package api

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type UserSession struct {
	UserID int64
}

type ourCustomKey string

const sessionKey ourCustomKey = "unique-session-key-for-our-example"

var (
	cookieStore = sessions.NewCookieStore([]byte("forDemo"))
)

func init() {
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func UserFromSession(req *http.Request) (UserSession, bool) {
	session, ok := req.Context().Value(sessionKey).(UserSession)
	if session.UserID < 1 {
		// Shouldnt happen
		return UserSession{}, false
	}
	return session, ok
}
