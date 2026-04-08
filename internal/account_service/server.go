package account_service

import (
	"TransactionManager/internal/account_service/contracts"
	"TransactionManager/internal/account_service/service"
	"TransactionManager/packages/public_response"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AccountHandlerServer is the HTTP layer.
type AccountHandlerServer struct {
	service *service.AccountService
}

// NewAccountHandlerServer creates a new handler and registers its routes.
func NewAccountHandlerServer(router chi.Router, svc *service.AccountService) *AccountHandlerServer {
	s := &AccountHandlerServer{service: svc}
	RegisterRoutes(router, s)
	return s
}

func (h *AccountHandlerServer) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req contracts.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		public_response.ToError(w, public_response.ErrValidation)
		return
	}
	if err := service.ValidateCreateAccountRequest(req); err != nil {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	response, err := h.service.CreateAccount(ctx, req)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.Created(w, response)
}

func (h *AccountHandlerServer) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountIDStr := chi.URLParam(r, "accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil || accountID <= 0 {
		public_response.ToErrorResponse(w, http.StatusBadRequest, "validation_failed", "invalid account_id")
		return
	}

	response, err := h.service.GetAccountByID(ctx, accountID)
	if err != nil {
		public_response.ToError(w, err)
		return
	}
	public_response.OK(w, response)
}
