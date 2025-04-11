package transaction

import (
	"context"
	db2 "github.com/Str1m/auth/internal/client/db"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// TODO: Проверить работу транзакций

type TxManager struct {
	db *db2.Client
}

func NewTransactionManager(db *db2.Client) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) ReadCommitted(ctx context.Context, f db2.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *TxManager) transaction(ctx context.Context, opts pgx.TxOptions, f db2.Handler) (err error) {
	tx, ok := ctx.Value(db2.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	tx, err = m.db.GetDB().BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctx = db2.MakeContextTx(ctx, tx)

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
