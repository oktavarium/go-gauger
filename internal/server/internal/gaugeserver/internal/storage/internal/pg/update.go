package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *storage) UpdateCounter(ctx context.Context, name string, val int64) (int64, error) {
	tx, err := s.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var newValue int64
	row := tx.QueryRow(ctx, `
		INSERT INTO counter (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO
		UPDATE SET value = counter.value + $2
		RETURNING value`, name, val)

	err = row.Scan(&newValue)
	if err != nil {
		return 0, fmt.Errorf("error occured on updating counter: %w", err)
	}

	return newValue, tx.Commit(ctx)
}
