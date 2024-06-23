package Workout

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
)

type UserSession struct {
	UserID int64
}

type ourCustomKey string

const sessionKey ourCustomKey = "unique-session-key-for-our-example"

func UserFromSession(req *http.Request) (UserSession, bool) {
	session, ok := req.Context().Value(sessionKey).(UserSession)
	if session.UserID < 1 {
		// Shouldnt happen
		return UserSession{}, false
	}
	return session, ok
}

func HandlecreateNewWorkout(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, ok := UserFromSession(req)
		if !ok {
			api.JSONError(wr, http.StatusForbidden, "Bad context")
			return
		}
		querier := store.New(db)

		res, err := querier.CreateUserWorkout(req.Context(), userDetails.UserID)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}

		json.NewEncoder(wr).Encode(&res)

	})
}
