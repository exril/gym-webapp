package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	Workout "github.com/extractings/gym-webapp/Workouts"
	controllers "github.com/extractings/gym-webapp/controllers/Users"
	"github.com/extractings/gym-webapp/internal"
	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/middlewares"
	"github.com/extractings/gym-webapp/store"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		internal.GetAsString("DB_USER", "local"),
		internal.GetAsString("DB_PASSWORD", "asecurepassword"),
		internal.GetAsString("DB_HOST", "localhost"),
		internal.GetAsInt("DB_PORT", 5432),
		internal.GetAsString("DB_NAME", "fullstackdb"),
	)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	createUserInDb(db)

	server := api.NewServer(internal.GetAsInt("SERVER_PORT", 9002))

	server.MustStart()
	defer server.Stop()

	defaultMiddleware := []mux.MiddlewareFunc{
		middlewares.JSONMiddleware,
		middlewares.CORSMiddleware(internal.GetAsSlice("CORS_WHITELIST",
			[]string{
				"http://localhost:9000",
				"http://0.0.0.0:9000",
			}, ","),
		),
	}

	// Handlers
	server.AddRoute("/login", controllers.HandleLogin(db), http.MethodPost, defaultMiddleware...)
	server.AddRoute("/logout", controllers.HandleLogout(), http.MethodGet, defaultMiddleware...)

	protectedMiddleware := append(defaultMiddleware, validCookieMiddleware(db))
	server.AddRoute("/checkSecret", checkSecret(db), http.MethodGet, protectedMiddleware...)

	// Workouts
	server.AddRoute("/workout", Workout.HandlecreateNewWorkout(db), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/workout", Workout.HandleListWorkouts(db), http.MethodGet, protectedMiddleware...)
	server.AddRoute("/workout/{workout_id}", Workout.HandleDeleteWorkout(db), http.MethodDelete, protectedMiddleware...)
	server.AddRoute("/workout/{workout_id}", Workout.HandleAddSet(db), http.MethodPost, protectedMiddleware...)
	// will cook in update workout function later (cannot figure out somethings)
	// server.AddRoute("/workout/{workout_id}/{set_id}", handleUpdateSet(db), http.MethodPut, protectedMiddleware...)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// Once received, we exit and the server is cleaned up
	<-sigChan
}

func createUserInDb(db *sql.DB) {

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
