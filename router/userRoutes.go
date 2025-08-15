package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/handlers"
	"github.com/sajidzamanme/emi-tracker/middlewares"
)

func InitUserRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"GET /users",
		manager.With(
			http.HandlerFunc(handlers.GetAllUsers),
		),
	)

	mux.Handle(
		"GET /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.GetUserByID),
		),
	)

	mux.Handle(
		"POST /users/signup",
		manager.With(
			http.HandlerFunc(handlers.InsertUser),
		),
	)

	mux.Handle(
		"POST /users/login",
		manager.With(
			http.HandlerFunc(handlers.UserLogin),
		),
	)

	mux.Handle(
		"PUT /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.UpdateUser),
		),
	)

	mux.Handle(
		"DELETE /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.DeleteUser),
		),
	)

	mux.Handle(
		"GET /users/{userID}/emirecords",
		manager.With(
			http.HandlerFunc(handlers.GetAllRecordsByUserID),
		),
	)
}
