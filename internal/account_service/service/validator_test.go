package service_test

import (
	"TransactionManager/internal/account_service/contracts"
	"TransactionManager/internal/account_service/service"
	"testing"
)

func TestValidateCreateAccountRequest(t *testing.T) {
	if err := service.ValidateCreateAccountRequest(contracts.CreateAccountRequest{DocumentNumber: ""}); err == nil {
		t.Fatalf("expected validation error for empty document number")
	}

	if err := service.ValidateCreateAccountRequest(contracts.CreateAccountRequest{DocumentNumber: "123"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
