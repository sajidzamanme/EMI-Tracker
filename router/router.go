package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/middlewares"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()

	manager := middlewares.NewManager()
	manager.Use(
		middlewares.Logger,
		middlewares.HandleCORS,
		middlewares.HandlePreflight,
	)
	wrappedMux := manager.WrapMux(mux)

	InitUserRoutes(mux, manager)
	InitEmiRecordRoutes(mux, manager)

	return wrappedMux
}
