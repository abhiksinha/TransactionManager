package repo

import (
	"TransactionManager/internal/account_service/model"

	"gorm.io/gorm"
)

// Repository provides access to the database for the account service.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ExecTxn executes the given function within a database transaction.
func (r *Repository) ExecTxn(fn func(repo *Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txnRepo := NewRepository(tx)
		return fn(txnRepo)
	})
}

// CreateAccount inserts a new account.
func (r *Repository) CreateAccount(account *model.Account) error {
	return r.db.Create(account).Error
}

// GetByID fetches an account by ID.
func (r *Repository) GetByID(id int64) (*model.Account, error) {
	var account model.Account
	err := r.db.First(&account, "id = ?", id).Error
	return &account, err
}

// GetByDocumentNumber fetches an account by document number.
func (r *Repository) GetByDocumentNumber(doc string) (*model.Account, error) {
	var account model.Account
	err := r.db.First(&account, "document_number = ?", doc).Error
	return &account, err
}
