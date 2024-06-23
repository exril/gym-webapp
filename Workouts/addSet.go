package Workout

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
	"github.com/gorilla/mux"
)

func HandleAddSet(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		workoutID, err := strconv.Atoi(mux.Vars(req)["workout_id"])
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad workout_id")
			return
		}

		type newSetRequest struct {
			ExerciseName string `json:"exercise_name,omitempty"`
			Weight       int    `json:"weight,omitempty"`
		}

		payload := newSetRequest{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		querier := store.New(db)

		set, err := querier.CreateDefaultSetForExercise(req.Context(),
			store.CreateDefaultSetForExerciseParams{
				WorkoutID:    int64(workoutID),
				ExerciseName: payload.ExerciseName,
				Weight:       int32(payload.Weight),
			})
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&set)
	})
}
