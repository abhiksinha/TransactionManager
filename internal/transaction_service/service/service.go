package service

import (
	accountmodel "TransactionManager/internal/account_service/model"
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/internal/transaction_service/model"
	"TransactionManager/internal/transaction_service/repo"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/public_response"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AccountReader provides account lookup for validation.
type AccountReader interface {
	GetByID(id int64) (*accountmodel.Account, error)
}

// TransactionService encapsulates the business logic for transactions.
type TransactionService struct {
	repo        *repo.Repository
	accountRepo AccountReader
	logger      *zap.Logger
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(repo *repo.Repository, accountRepo AccountReader, logger *zap.Logger) *TransactionService {
	return &TransactionService{repo: repo, accountRepo: accountRepo, logger: logger}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req contracts.CreateTransactionRequest) (*contracts.TransactionResponse, error) {
	log := logger.FromContext(ctx, s.logger)

	if _, err := s.accountRepo.GetByID(req.AccountID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Account not found", zap.Int64("account_id", req.AccountID))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching account", zap.Error(err))
		return nil, err
	}

	opType, err := s.repo.GetOperationTypeByID(req.OperationTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Operation type not found", zap.Int64("operation_type_id", req.OperationTypeID))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching operation type", zap.Error(err))
		return nil, err
	}

	amountMinor, err := amountToMinorUnits(req.Amount)
	if err != nil {
		log.Warn("Invalid transaction amount", zap.Error(err))
		return nil, public_response.ErrValidation
	}

	if opType.TransactionType != transactionTypeDebit && opType.TransactionType != transactionTypeCredit {
		log.Warn("Invalid transaction type on operation type", zap.String("transaction_type", opType.TransactionType))
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
		log.Error("Error creating transaction", zap.Error(err))
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
