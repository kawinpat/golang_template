package routes

import (
	"golang_template/handler"
	"golang_template/handler/auth"
	"golang_template/handler/users"
	"golang_template/helper"

	"github.com/gorilla/mux"
)

// InitialRoutes sets up the routes and returns the router
func InitialRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health Check
	r.HandleFunc("/", handler.HealthCheck).Methods("GET")

	// signin
	r.HandleFunc("/signin", auth.SignIn).Methods("POST")
	r.HandleFunc("/signout", auth.SignOut).Methods("POST")

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.Use(helper.TokenValidationMiddleware) // Check exp token middleware
	authRoutes.HandleFunc("/sessions/cleanup", auth.CleanupExpiredSessions).Methods("POST")

	userRoutes := authRoutes.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/", users.HealthCheck).Methods("GET")
	userRoutes.HandleFunc("/list", users.List).Methods("GET")
	userRoutes.HandleFunc("/detail/{id}", users.Detail).Methods("GET")
	userRoutes.HandleFunc("/create", users.Create).Methods("POST")
	userRoutes.HandleFunc("/update/{id}", users.Update).Methods("put")

	return r
}
