package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/extractings/gym-webapp/internal"
	"github.com/extractings/gym-webapp/store"
	"github.com/lib/pq"
)

func CreateUserInDb(db *sql.DB) {

	ctx := context.Background()
	querier := store.New(db)

	log.Println("Creating user@user...")
	hashPwd := internal.HashPassword("password")

	_, err := querier.CreateUsers(ctx, store.CreateUsersParams{
		UserName:     "user@user",
		PasswordHash: hashPwd,
		Name:         "Dummy user",
	})

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		log.Println("Dummy User already present")
		return
	}

	if err != nil {
		log.Println("Failed to create user:", err)
	}
}
