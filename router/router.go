package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/handlers"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// User Handlers
	mux.HandleFunc("GET /users", handlers.GetAllUsersHandler)
	mux.HandleFunc("GET /users/{userID}", handlers.GetUserByIDHandler)
	mux.HandleFunc("POST /users/signup", handlers.InsertUserHandler)
	mux.HandleFunc("POST /users/login", handlers.UserLoginHandler)
	mux.HandleFunc("PUT /users/{userID}", handlers.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/{userID}", handlers.DeleteUserHandler)
	mux.HandleFunc("GET /users/{userID}/emirecords", handlers.GetAllRecordsByUserIDHandler)

	// EMIRecord Handlers
	mux.HandleFunc("GET /emirecords/{recordID}", handlers.GetRecordByRecordID)
	mux.HandleFunc("POST /emirecords/{userID}", handlers.InsertRecordByUserID)
	mux.HandleFunc("PUT /emirecords/{recordID}", handlers.UpdateRecordByRecordID)
	mux.HandleFunc("DELETE /emirecords/{recordID}", handlers.DeleteRecordByRecordID)
	mux.HandleFunc("GET /emirecords/{recordID}/payinstallment", handlers.PayInstallment)

	return mux
}
