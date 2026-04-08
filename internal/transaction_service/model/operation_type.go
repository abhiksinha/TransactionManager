package model

// OperationType represents the operation_types table.
type OperationType struct {
	ID              int64  `gorm:"primaryKey"`
	Description     string `gorm:"type:varchar(255);not null"`
	TransactionType string `gorm:"type:varchar(10);not null"`
}

// TableName explicitly sets the table name for the OperationType model.
func (o *OperationType) TableName() string {
	return "operation_types"
}
