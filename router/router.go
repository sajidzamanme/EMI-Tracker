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
	mux.HandleFunc("GET /users/{userID}/emirecords", handlers.GetAllRecordsByUserID)

	// EMIRecord Handlers
	mux.HandleFunc("GET /emirecords/{recordID}", handlers.GetRecordByRecordID)
	mux.HandleFunc("POST /emirecords/{userID}", handlers.PostRecordByUserID)
	mux.HandleFunc("PUT /emirecords/{recordID}", handlers.PutRecordByRecordID)
	mux.HandleFunc("DELETE /emirecords/{recordID}", handlers.DeleteRecordByRecordID)

	return mux
}
