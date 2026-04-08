package model

import "time"

// Account represents the accounts table in the database.
type Account struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	DocumentNumber string    `gorm:"type:varchar(20);unique;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

// TableName explicitly sets the table name for the Account model.
func (a *Account) TableName() string {
	return "accounts"
}
