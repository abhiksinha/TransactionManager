package account_service_test

import (
	account_service "TransactionManager/internal/account_service"
	"TransactionManager/internal/account_service/contracts"
	"TransactionManager/packages/public_response"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeAccountService struct {
	createResp *contracts.AccountResponse
	getResp    *contracts.AccountResponse
	err        error
}

func (f *fakeAccountService) CreateAccount(ctx context.Context, req contracts.CreateAccountRequest) (*contracts.AccountResponse, error) {
	return f.createResp, f.err
}

func (f *fakeAccountService) GetAccountByID(ctx context.Context, id int64) (*contracts.AccountResponse, error) {
	return f.getResp, f.err
}

func TestAccountServerCreateAccountInvalidJSON(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString("{invalid"))
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

func TestAccountServerCreateAccountValidation(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":""}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAccountServerCreateAccountSuccess(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{
		createResp: &contracts.AccountResponse{AccountID: 10, DocumentNumber: "123"},
	}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":"123"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestAccountServerGetAccountInvalidID(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAccountServerGetAccountNotFound(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{err: public_response.ErrNotFound}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestAccountServerGetAccountSuccess(t *testing.T) {
	router := chi.NewRouter()
	svc := &fakeAccountService{getResp: &contracts.AccountResponse{AccountID: 1, DocumentNumber: "123"}}
	account_service.NewAccountHandlerServer(router, svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
