package main

import (
	"context"
	"database/sql"
	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
	"net/http"
)


type UserSession struct {
	UserID int64
}

type ourCustomKey string

const sessionKey ourCustomKey = "unique-session-key-for-our-example"

func validCookieMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			session, err := cookieStore.Get(req, "session-name")
			if err != nil {
				api.JSONError(wr, http.StatusInternalServerError, "Session Error")
				return
			}

			userID, userIDOK := session.Values["userID"].(int64)
			isAuthd, isAuthdOK := session.Values["userAuthenticated"].(bool)

			// We could put with the above but lets keep our logic simple
			if !userIDOK || !isAuthdOK {
				api.JSONError(wr, http.StatusInternalServerError, "Session Error")
				return
			}

			if !isAuthd || userID < 1 {
				api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
				return
			}

			querier := store.New(db)
			user, err := querier.GetUser(req.Context(), int64(userID))
			if err != nil || user.UserID < 1 {
				api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
				return
			}

			ctx := context.WithValue(req.Context(), sessionKey, UserSession{
				UserID: user.UserID,
			})
			h.ServeHTTP(wr, req.WithContext(ctx))
		})
	}
}

func userFromSession(req *http.Request) (UserSession, bool) {
	session, ok := req.Context().Value(sessionKey).(UserSession)
	if session.UserID < 1 {
		// Shouldnt happen
		return UserSession{}, false
	}
	return session, ok
}
