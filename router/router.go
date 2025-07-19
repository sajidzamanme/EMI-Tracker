package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/handlers"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// User Handlers
	mux.HandleFunc("GET /users", handlers.GetAllUsers)
	mux.HandleFunc("GET /users/{userID}", handlers.GetUserByID)
	mux.HandleFunc("POST /users", handlers.PostUser)
	mux.HandleFunc("PUT /users/{userID}", handlers.PutUser)
	mux.HandleFunc("DELETE /users/{userID}", handlers.DeleteUser)
	mux.HandleFunc("GET /users/{userID}/subscriptions", handlers.GetAllSubsByUserID)

	// Subscription Handlers
	mux.HandleFunc("GET /subscriptions/{subID}", handlers.GetSubByID)
	mux.HandleFunc("POST /subscriptions/{userID}", handlers.PostSubByUserID)
	mux.HandleFunc("PUT /subscriptions/{subID}", handlers.PutSubBySubID)
	mux.HandleFunc("DELETE /subscriptions/{subID}", handlers.DeleteSubBySubID)

	return mux
}
