package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/extractings/gym-webapp/internal"
	"github.com/extractings/gym-webapp/store"
	"github.com/lib/pq"
)

func initDatabase() {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		GetAsString("DB_USER", "postgres"),
		GetAsString("DB_PASSWORD", "mysecretpassword"),
		GetAsString("DB_HOST", "localhost"),
		GetAsInt("DB_PORT", 5432),
		GetAsString("DB_NAME", "postgres"),
	)

	// Open the database
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// Connectivity check
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	// Create the store
	dbQuery = store.New(db)

	ctx := context.Background()

	CreateUserInDb(db)

	if err != nil {
		os.Exit(1)
	}
}