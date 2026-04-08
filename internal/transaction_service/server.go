package transaction_service

import (
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/internal/transaction_service/service"
	"TransactionManager/packages/public_response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TransactionHandlerServer is the HTTP layer.
type TransactionHandlerServer struct {
	service *service.TransactionService
}

// NewTransactionHandlerServer creates a new handler and registers its routes.
func NewTransactionHandlerServer(router chi.Router, svc *service.TransactionService) *TransactionHandlerServer {
	s := &TransactionHandlerServer{service: svc}
	RegisterRoutes(router, s)
	return s
}

func (h *TransactionHandlerServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateTransactionRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	response, err := h.service.CreateTransaction(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}
