package transaction_service_test

import (
	transaction_service "TransactionManager/internal/transaction_service"
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/packages/public_response"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeTransactionService struct {
	resp *contracts.TransactionResponse
	err  error
}

func (f *fakeTransactionService) CreateTransaction(ctx context.Context, req contracts.CreateTransactionRequest) (*contracts.TransactionResponse, error) {
	return f.resp, f.err
}

func TestTransactionServerCreateTransactionInvalidJSON(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeTransactionService{}
	transaction_service.NewTransactionHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString("{invalid"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var resp public_response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Code != "validation_failed" {
		t.Fatalf("unexpected error code: %s", resp.Code)
	}
}

func TestTransactionServerCreateTransactionValidation(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeTransactionService{}
	transaction_service.NewTransactionHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":0}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestTransactionServerCreateTransactionNotFound(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeTransactionService{err: public_response.ErrNotFound}
	transaction_service.NewTransactionHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":1,"operation_type_id":1,"amount":10}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestTransactionServerCreateTransactionSuccess(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeTransactionService{
		resp: &contracts.TransactionResponse{TransactionID: 1, AccountID: 1, OperationTypeID: 2, Amount: 10},
	}
	transaction_service.NewTransactionHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":1,"operation_type_id":2,"amount":10}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}
