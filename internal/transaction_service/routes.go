package transaction_service

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers transaction-related routes on the provided router.
func RegisterRoutes(router chi.Router, handler *TransactionHandlerServer) {
	router.Route("/transactions", func(r chi.Router) {
		r.Post("/", handler.CreateTransaction)
	})
}
