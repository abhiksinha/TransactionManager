package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTransactionsTable, downTransactionsTable)
}

func upTransactionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			account_id BIGINT NOT NULL REFERENCES accounts(id),
			operation_type_id BIGINT NOT NULL REFERENCES operation_types(id),
			amount BIGINT NOT NULL,
			event_date TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	return err
}

func downTransactionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS transactions;`)
	return err
}
