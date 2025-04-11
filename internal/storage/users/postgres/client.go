package postgres

import (
	"context"
	"github.com/Str1m/auth/internal/client/db"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type key string

const (
	TxKey key = "tx"
)

type ClientPG struct {
	pg *pgxpool.Pool
}

func NewClientPG(pg *pgxpool.Pool) *ClientPG {
	return &ClientPG{pg: pg}
}

func (p *ClientPG) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}
	return p.pg.Exec(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}
	return p.pg.Query(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.pg.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) Ping(ctx context.Context) error {
	return p.pg.Ping(ctx)
}

func (p *ClientPG) Close() {
	p.pg.Close()
}

func (p *ClientPG) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, row)
}

func (p *ClientPG) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, row)
}
