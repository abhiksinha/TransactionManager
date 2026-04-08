package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOperationTypesTable, downOperationTypesTable)
}

func upOperationTypesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS operation_types (
			id BIGINT PRIMARY KEY,
			description VARCHAR(255) NOT NULL,
			transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('debit', 'credit'))
		);
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO operation_types (id, description, transaction_type) VALUES
			(1, 'Normal Purchase', 'debit'),
			(2, 'Purchase with installments', 'debit'),
			(3, 'Withdrawal', 'debit'),
			(4, 'Credit Voucher', 'credit')
		ON CONFLICT (id) DO NOTHING;
	`)
	return err
}

func downOperationTypesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS operation_types;`)
	return err
}
