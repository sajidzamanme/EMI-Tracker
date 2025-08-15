package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (m *Manager) Use(middlewares ...Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
}

func (m * Manager) With(handler http.Handler, middlewares ...Middleware) http.Handler {
	h := handler

	for _, mw := range middlewares {
		h = mw(h)
	}

	for _, mw := range m.globalMiddlewares {
		h = mw(h)
	}

	return h
}