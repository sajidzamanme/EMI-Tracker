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
			http.HandlerFunc(handlers.GetAllUsersHandler),
		),
	)

	mux.Handle(
		"GET /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.GetUserByIDHandler),
		),
	)

	mux.Handle(
		"POST /users/signup",
		manager.With(
			http.HandlerFunc(handlers.InsertUserHandler),
		),
	)

	mux.Handle(
		"POST /users/login",
		manager.With(
			http.HandlerFunc(handlers.UserLoginHandler),
		),
	)

	mux.Handle(
		"PUT /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.UpdateUserHandler),
		),
	)

	mux.Handle(
		"DELETE /users/{userID}",
		manager.With(
			http.HandlerFunc(handlers.DeleteUserHandler),
		),
	)

	mux.Handle(
		"GET /users/{userID}/emirecords",
		manager.With(
			http.HandlerFunc(handlers.GetAllRecordsByUserIDHandler),
		),
	)
}
