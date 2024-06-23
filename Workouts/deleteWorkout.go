package Workout

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
)

func HandleListWorkouts(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, ok := UserFromSession(req)
		if !ok {
			api.JSONError(wr, http.StatusForbidden, "Bad context")
			return
		}

		querier := store.New(db)
		workouts, err := querier.GetWorkoutsForUserID(req.Context(), userDetails.UserID)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&workouts)
	})
}
