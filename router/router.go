package router

import (
	"net/http"

	"github.com/sajidzamanme/emi-tracker/middlewares"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	manager := middlewares.NewManager()
	manager.Use(middlewares.Logger)

	InitUserRoutes(mux, manager)
	InitEmiRecordRoutes(mux, manager)

	return mux
}
