package controllers

import (
	"net/http"

	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/extractings/gym-webapp/internal"
	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
	"github.com/goccy/go-json"
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

func HandleLogin(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		type loginRequest struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		payload := loginRequest{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		querier := store.New(db)
		user, err := querier.GetUserByName(req.Context(), payload.Username)
		if errors.Is(err, sql.ErrNoRows) || !internal.CheckPasswordHash(payload.Password, user.PasswordHash) {
			api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
			return
		}
		if err != nil {
			log.Println("Received error looking up user", err)
			api.JSONError(wr, http.StatusInternalServerError, "Couldn't log you in due to a server error")
			return
		}

		session, err := cookieStore.Get(req, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
		}
		session.Values["userAuthenticated"] = true
		session.Values["userID"] = user.UserID
		session.Save(req, wr)
	})
}

func CheckSecret(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, _ := UserFromSession(req)

		querier := store.New(db)
		user, err := querier.GetUser(req.Context(), userDetails.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			api.JSONError(wr, http.StatusForbidden, "User not found")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Hello there %s", user.Name))
	})
}

func HandleLogout() http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		session, err := cookieStore.Get(req, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
			return
		}

		session.Options.MaxAge = -1 // deletes
		session.Values["userID"] = int64(-1)
		session.Values["userAuthenticated"] = false

		err = session.Save(req, wr)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
			return
		}

		api.JSONMessage(wr, http.StatusOK, "logout successful")
	})
}
