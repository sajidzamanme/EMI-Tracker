package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/handlers"
	"github.com/sajidzamanme/emi-tracker/middlewares"
)

func InitEmiRecordRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"GET /emirecords/{recordID}",
		manager.With(
			http.HandlerFunc(handlers.GetRecordByRecordIDHandler),
		),
	)

	mux.Handle(
		"POST /emirecords/{userID}",
		manager.With(
			http.HandlerFunc(handlers.InsertRecordByUserIDHandler),
		),
	)

	mux.Handle(
		"PUT /emirecords/{recordID}",
		manager.With(
			http.HandlerFunc(handlers.UpdateRecordByRecordIDHandler),
		),
	)

	mux.Handle(
		"DELETE /emirecords/{recordID}",
		manager.With(
			http.HandlerFunc(handlers.DeleteRecordByRecordIDHandler),
		),
	)

	mux.Handle(
		"GET /emirecords/{recordID}/payinstallment",
		manager.With(
			http.HandlerFunc(handlers.PayInstallmentHandler),
		),
	)
}
