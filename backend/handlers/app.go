package handlers

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	Workout "github.com/extractings/gym-webapp/Workouts"
	controllers "github.com/extractings/gym-webapp/controllers/Users"
	"github.com/extractings/gym-webapp/internal"
	"github.com/extractings/gym-webapp/internal/api"
	"github.com/extractings/gym-webapp/middlewares"
	"github.com/gorilla/mux"
)

func LoadApplication() {
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

	// our session protected middlewares
	protectedMiddleware := append(defaultMiddleware, ValidCookieMiddleware(db))
	server.AddRoute("/checkSecret", controllers.CheckSecret(db), http.MethodGet, protectedMiddleware...)

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
	server.Use(middlewares.BasicMiddleware())

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
