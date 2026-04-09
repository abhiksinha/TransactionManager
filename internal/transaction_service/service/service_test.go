package service_test

import (
	accountcontracts "TransactionManager/internal/account_service/contracts"
	accountservice "TransactionManager/internal/account_service/service"
	"TransactionManager/internal/transaction_service/contracts"
	"TransactionManager/internal/transaction_service/model"
	"TransactionManager/internal/transaction_service/repo"
	"TransactionManager/internal/transaction_service/service"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/public_response"
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeAccountService struct {
	account *accountcontracts.AccountResponse
	err     error
}

func (f *fakeAccountService) CreateAccount(ctx context.Context, req accountcontracts.CreateAccountRequest) (*accountcontracts.AccountResponse, error) {
	return f.account, f.err
}

func (f *fakeAccountService) GetAccountByID(ctx context.Context, id int64) (*accountcontracts.AccountResponse, error) {
	return f.account, f.err
}

func newTxnTestDB(t *testing.T, name string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.OperationType{}, &model.Transaction{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newTxnLogger(t *testing.T) *logger.Logger {
	t.Helper()
	l, err := logger.New(false)
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	return l
}

func TestTransactionServiceCreateTransaction(t *testing.T) {
	db := newTxnTestDB(t, t.Name())
	if err := db.Create(&model.OperationType{ID: 4, Description: "PAYMENT", TransactionType: "credit"}).Error; err != nil {
		t.Fatalf("seed op type: %v", err)
	}

	accountSvc := &fakeAccountService{
		account: &accountcontracts.AccountResponse{AccountID: 1, DocumentNumber: "123"},
	}
	svc := service.NewTransactionService(repo.NewRepository(db), accountSvc, newTxnLogger(t))

	resp, err := svc.CreateTransaction(context.Background(), contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          12.34,
	})
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}
	if resp.TransactionID == 0 {
		t.Fatalf("expected transaction id to be set")
	}
	if resp.Amount != 12.34 {
		t.Fatalf("unexpected signed amount: %v", resp.Amount)
	}
}

func TestTransactionServiceCreateTransactionDebitStoresNegative(t *testing.T) {
	db := newTxnTestDB(t, t.Name())
	if err := db.Create(&model.OperationType{ID: 1, Description: "Normal Purchase", TransactionType: "debit"}).Error; err != nil {
		t.Fatalf("seed op type: %v", err)
	}

	accountSvc := &fakeAccountService{
		account: &accountcontracts.AccountResponse{AccountID: 1, DocumentNumber: "123"},
	}
	svc := service.NewTransactionService(repo.NewRepository(db), accountSvc, newTxnLogger(t))

	resp, err := svc.CreateTransaction(context.Background(), contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          10.00,
	})
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}
	if resp.Amount != -10.00 {
		t.Fatalf("expected signed response amount -10.00, got %v", resp.Amount)
	}

	var txn model.Transaction
	if err := db.First(&txn, "id = ?", resp.TransactionID).Error; err != nil {
		t.Fatalf("fetch transaction: %v", err)
	}
	if txn.Amount != -1000 {
		t.Fatalf("expected stored amount -1000, got %d", txn.Amount)
	}
}

func TestTransactionServiceCreateTransactionInvalidAmount(t *testing.T) {
	db := newTxnTestDB(t, t.Name())
	if err := db.Create(&model.OperationType{ID: 4, Description: "PAYMENT", TransactionType: "credit"}).Error; err != nil {
		t.Fatalf("seed op type: %v", err)
	}

	accountSvc := &fakeAccountService{
		account: &accountcontracts.AccountResponse{AccountID: 1, DocumentNumber: "123"},
	}
	svc := service.NewTransactionService(repo.NewRepository(db), accountSvc, newTxnLogger(t))

	_, err := svc.CreateTransaction(context.Background(), contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          12.345,
	})
	if err == nil || err != public_response.ErrValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTransactionServiceCreateTransactionAccountNotFound(t *testing.T) {
	db := newTxnTestDB(t, t.Name())
	if err := db.Create(&model.OperationType{ID: 4, Description: "PAYMENT", TransactionType: "credit"}).Error; err != nil {
		t.Fatalf("seed op type: %v", err)
	}

	accountSvc := &fakeAccountService{err: public_response.ErrNotFound}
	svc := service.NewTransactionService(repo.NewRepository(db), accountSvc, newTxnLogger(t))

	_, err := svc.CreateTransaction(context.Background(), contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          12.34,
	})
	if err == nil || err != public_response.ErrNotFound {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestTransactionServiceCreateTransactionOperationTypeNotFound(t *testing.T) {
	db := newTxnTestDB(t, t.Name())

	accountSvc := &fakeAccountService{
		account: &accountcontracts.AccountResponse{AccountID: 1, DocumentNumber: "123"},
	}
	svc := service.NewTransactionService(repo.NewRepository(db), accountSvc, newTxnLogger(t))

	_, err := svc.CreateTransaction(context.Background(), contracts.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 999,
		Amount:          12.34,
	})
	if err == nil || err != public_response.ErrNotFound {
		t.Fatalf("expected not found error, got %v", err)
	}
}

var _ accountservice.AccountService = (*fakeAccountService)(nil)
