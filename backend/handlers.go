package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/store"
)

func checkSecret(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, _ := api.UserFromSession(req)

		querier := store.New(db)
		user, err := querier.GetUser(req.Context(), userDetails.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			api.JSONError(wr, http.StatusForbidden, "User not found")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Hello there %s", user.Name))
	})
}
