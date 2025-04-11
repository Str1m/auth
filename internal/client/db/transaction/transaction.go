package transaction

import (
	"context"
	"github.com/Str1m/auth/internal/client/db"
	"github.com/Str1m/auth/internal/client/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, f db.Handler) (err error) {
	tx, ok := ctx.Value(postgres.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctx = postgres.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); err != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)

			}
			return
		}

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "commit failed")
			}
		}
	}()

	if err = f(ctx); err != nil {
		err = errors.Wrap(err, "failed exec code inside transaction")
	}
	return nil
}
