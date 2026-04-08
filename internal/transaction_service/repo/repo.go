package repo

import (
	"TransactionManager/internal/transaction_service/model"

	"gorm.io/gorm"
)

// Repository provides access to the database for the transaction service.
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

// GetOperationTypeByID fetches an operation type by ID.
func (r *Repository) GetOperationTypeByID(id int64) (*model.OperationType, error) {
	var op model.OperationType
	err := r.db.First(&op, "id = ?", id).Error
	return &op, err
}

// CreateTransaction inserts a new transaction.
func (r *Repository) CreateTransaction(txn *model.Transaction) error {
	return r.db.Create(txn).Error
}
