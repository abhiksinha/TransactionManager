package service

import (
	"TransactionManager/internal/transaction_service/contracts"
	"errors"
	"math"
)

// ValidateCreateTransactionRequest validates the create transaction request.
func ValidateCreateTransactionRequest(req contracts.CreateTransactionRequest) error {
	if req.AccountID <= 0 {
		return errors.New("account_id is required")
	}
	if req.OperationTypeID <= 0 {
		return errors.New("operation_type_id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	minor := req.Amount * 100
	if math.Abs(minor-math.Round(minor)) > 0.000001 {
		return errors.New("amount must have at most two decimal places")
	}
	return nil
}
