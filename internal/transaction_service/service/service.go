package service

import (
	accountservice "TransactionManager/internal/account_service/service"
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/internal/transaction_service/model"
	"TransactionManager/internal/transaction_service/repo"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/public_response"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransactionService defines the behavior for transaction operations.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req contracts.CreateTransactionRequest) (*contracts.TransactionResponse, error)
}

// transactionService encapsulates the business logic for transactions.
type transactionService struct {
	repo       *repo.Repository
	accountSvc accountservice.AccountService
	logger     *logger.Logger
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(repo *repo.Repository, accountSvc accountservice.AccountService, logger *logger.Logger) TransactionService {
	return &transactionService{repo: repo, accountSvc: accountSvc, logger: logger}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req contracts.CreateTransactionRequest) (*contracts.TransactionResponse, error) {
	opType, err := s.verify(ctx, req)
	if err != nil {
		return nil, err
	}

	amountMinor, err := amountToMinorUnits(req.Amount)
	if err != nil {
		s.logger.Warn(ctx, "Invalid transaction amount", zap.Error(err))
		return nil, public_response.ErrValidation
	}

	if opType.TransactionType != transactionTypeDebit && opType.TransactionType != transactionTypeCredit {
		s.logger.Warn(ctx, "Invalid transaction type on operation type", zap.String("transaction_type", opType.TransactionType))
		return nil, public_response.ErrValidation
	}

	txn := &model.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          amountMinor,
		EventDate:       nowIST(),
	}

	if err := s.repo.ExecTxn(func(txnRepo *repo.Repository) error {
		return s.repo.CreateTransaction(txn)
	}); err != nil {
		s.logger.Error(ctx, "Error creating transaction", zap.Error(err))
		return nil, err
	}

	return &contracts.TransactionResponse{
		TransactionID:   txn.ID,
		AccountID:       txn.AccountID,
		OperationTypeID: txn.OperationTypeID,
		Amount:          signedAmount(amountMinor, opType.TransactionType),
		EventDate:       txn.EventDate,
	}, nil
}

func (s *transactionService) verify(ctx context.Context, req contracts.CreateTransactionRequest) (*model.OperationType, error) {
	if _, err := s.accountSvc.GetAccountByID(ctx, req.AccountID); err != nil {
		if errors.Is(err, public_response.ErrNotFound) {
			s.logger.Warn(ctx, "Account not found", zap.Int64("account_id", req.AccountID))
			return nil, public_response.ErrNotFound
		}
		s.logger.Error(ctx, "Error fetching account", zap.Error(err))
		return nil, err
	}

	opType, err := s.repo.GetOperationTypeByID(req.OperationTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn(ctx, "Operation type not found", zap.Int64("operation_type_id", req.OperationTypeID))
			return nil, public_response.ErrNotFound
		}
		s.logger.Error(ctx, "Error fetching operation type", zap.Error(err))
		return nil, err
	}

	return opType, nil
}
