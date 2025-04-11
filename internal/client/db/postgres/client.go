package postgres

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler func(ctx context.Context) error

type Query struct {
	Name     string
	QueryRaw string
}
type key string

const (
	TxKey key = "tx"
)

type ClientPG struct {
	db *pgxpool.Pool
}

func NewClient(db *pgxpool.Pool) *ClientPG {
	return &ClientPG{db: db}
}

func (p *ClientPG) GetDB() *pgxpool.Pool {
	return p.db
}

func (p *ClientPG) ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}
	return p.db.Exec(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}
	return p.db.Query(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.db.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *ClientPG) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func (p *ClientPG) Close() {
	p.db.Close()
}

func (p *ClientPG) ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, row)
}

func (p *ClientPG) ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, row)
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
