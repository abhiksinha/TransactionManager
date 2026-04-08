package service_test

import (
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/internal/transaction_service/service"
	"testing"
)

func TestValidateCreateTransactionRequest(t *testing.T) {
	if err := service.ValidateCreateTransactionRequest(contracts.CreateTransactionRequest{}); err == nil {
		t.Fatalf("expected error for missing fields")
	}

	if err := service.ValidateCreateTransactionRequest(contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          0,
	}); err == nil {
		t.Fatalf("expected error for zero amount")
	}

	if err := service.ValidateCreateTransactionRequest(contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          12.345,
	}); err == nil {
		t.Fatalf("expected error for amount with >2 decimals")
	}

	if err := service.ValidateCreateTransactionRequest(contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          12.34,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
