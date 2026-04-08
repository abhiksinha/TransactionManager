package account_service

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers account-related routes on the provided router.
func RegisterRoutes(router chi.Router, handler *AccountHandlerServer) {
	router.Route("/accounts", func(r chi.Router) {
		r.Post("/", handler.CreateAccount)
		r.Get("/{accountId}", handler.GetAccount)
	})
}
