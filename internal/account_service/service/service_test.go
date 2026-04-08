package service_test

import (
	"TransactionManager/internal/account_service/contracts"
	"TransactionManager/internal/account_service/model"
	"TransactionManager/internal/account_service/repo"
	"TransactionManager/internal/account_service/service"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/public_response"
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T, name string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.Account{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newTestLogger(t *testing.T) *logger.Logger {
	t.Helper()
	l, err := logger.New(false)
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	return l
}

func TestAccountServiceCreateAccount(t *testing.T) {
	db := newTestDB(t, t.Name())
	svc := service.NewAccountService(repo.NewRepository(db), newTestLogger(t))

	resp, err := svc.CreateAccount(context.Background(), contracts.CreateAccountRequest{
		DocumentNumber: "12345678900",
	})
	if err != nil {
		t.Fatalf("create account: %v", err)
	}
	if resp.AccountID == 0 {
		t.Fatalf("expected account id to be set")
	}
	if resp.DocumentNumber != "12345678900" {
		t.Fatalf("unexpected document number: %s", resp.DocumentNumber)
	}
}

func TestAccountServiceCreateAccountDuplicate(t *testing.T) {
	db := newTestDB(t, t.Name())
	svc := service.NewAccountService(repo.NewRepository(db), newTestLogger(t))

	_, err := svc.CreateAccount(context.Background(), contracts.CreateAccountRequest{
		DocumentNumber: "12345678900",
	})
	if err != nil {
		t.Fatalf("create account: %v", err)
	}

	_, err = svc.CreateAccount(context.Background(), contracts.CreateAccountRequest{
		DocumentNumber: "12345678900",
	})
	if err == nil || err != public_response.ErrDuplicateEntry {
		t.Fatalf("expected duplicate entry error, got %v", err)
	}
}

func TestAccountServiceGetAccountByIDNotFound(t *testing.T) {
	db := newTestDB(t, t.Name())
	svc := service.NewAccountService(repo.NewRepository(db), newTestLogger(t))

	_, err := svc.GetAccountByID(context.Background(), 999)
	if err == nil || err != public_response.ErrNotFound {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestAccountServiceGetAccountByID(t *testing.T) {
	db := newTestDB(t, t.Name())
	svc := service.NewAccountService(repo.NewRepository(db), newTestLogger(t))

	account := &model.Account{DocumentNumber: "12345678900"}
	if err := db.Create(account).Error; err != nil {
		t.Fatalf("seed account: %v", err)
	}

	resp, err := svc.GetAccountByID(context.Background(), account.ID)
	if err != nil {
		t.Fatalf("get account: %v", err)
	}
	if resp.AccountID != account.ID {
		t.Fatalf("unexpected account id: %d", resp.AccountID)
	}
}
