package model

import "time"

// Transaction represents the transactions table in the database.
type Transaction struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	AccountID       int64     `gorm:"not null"`
	OperationTypeID int64     `gorm:"not null"`
	Amount          int64     `gorm:"not null"`
	EventDate       time.Time `gorm:"autoCreateTime"`
}

// TableName explicitly sets the table name for the Transaction model.
func (t *Transaction) TableName() string {
	return "transactions"
}
