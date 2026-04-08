package service

import (
	"TransactionManager/internal/account_service/contracts"
	"TransactionManager/internal/account_service/model"
	"TransactionManager/internal/account_service/repo"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/public_response"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AccountService defines the behavior for account operations.
type AccountService interface {
	CreateAccount(ctx context.Context, req contracts.CreateAccountRequest) (*contracts.AccountResponse, error)
	GetAccountByID(ctx context.Context, id int64) (*contracts.AccountResponse, error)
}

// accountService encapsulates the business logic for accounts.
type accountService struct {
	repo   *repo.Repository
	logger *zap.Logger
}

// NewAccountService creates a new AccountService.
func NewAccountService(repo *repo.Repository, logger *zap.Logger) AccountService {
	return &accountService{repo: repo, logger: logger}
}

func (s *accountService) CreateAccount(ctx context.Context, req contracts.CreateAccountRequest) (*contracts.AccountResponse, error) {
	log := logger.FromContext(ctx, s.logger)

	existing, err := s.repo.GetByDocumentNumber(req.DocumentNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Error checking for existing account", zap.Error(err))
		return nil, err
	}
	if err == nil && existing.ID != 0 {
		log.Warn("Attempted to create a duplicate account", zap.String("document_number", req.DocumentNumber))
		return nil, public_response.ErrDuplicateEntry
	}

	account := &model.Account{DocumentNumber: req.DocumentNumber}
	if err := s.repo.ExecTxn(func(txnRepo *repo.Repository) error {
		return s.repo.CreateAccount(account)
	}); err != nil {
		log.Error("Error during account creation transaction", zap.Error(err))
		return nil, err
	}

	return &contracts.AccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, id int64) (*contracts.AccountResponse, error) {
	log := logger.FromContext(ctx, s.logger)

	account, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Account not found", zap.Int64("id", id))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching account", zap.Error(err))
		return nil, err
	}

	return &contracts.AccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}, nil
}
