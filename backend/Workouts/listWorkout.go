package Workout

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
	"github.com/gorilla/mux"
)

func HandleDeleteWorkout(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, ok := UserFromSession(req)
		if !ok {
			api.JSONError(wr, http.StatusForbidden, "Bad context")
			return
		}

		workoutID, err := strconv.Atoi(mux.Vars(req)["workout_id"])
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad workout_id")
			return
		}

		err = store.New(db).DeleteWorkoutByIDForUser(req.Context(), store.DeleteWorkoutByIDForUserParams{
			UserID:    userDetails.UserID,
			WorkoutID: int64(workoutID),
		})

		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad workout_id")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Workout %d is deleted", workoutID))
	})
}
