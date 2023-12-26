package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var createTablesQuery string = `
CREATE TABLE IF NOT EXISTS
gauge
(
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE,
	value DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS
counter
(
	id SERIAL PRIMARY KEY,
 	name TEXT UNIQUE,
 	value BIGINT
);
`

func (s *storage) bootstrap(ctx context.Context) (err error) {
	tx, err := s.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}

	}()

	_, err = tx.Exec(ctx, createTablesQuery)
	if err != nil {
		return fmt.Errorf("error occured on creating tables: %w", err)
	}

	return nil
}
