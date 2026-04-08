package service

import (
	"TransactionManager/internal/account_service/contracts"
	"errors"
	"strings"
)

// ValidateCreateAccountRequest validates the create account request.
func ValidateCreateAccountRequest(req contracts.CreateAccountRequest) error {
	if strings.TrimSpace(req.DocumentNumber) == "" {
		return errors.New("document_number is required")
	}
	return nil
}
