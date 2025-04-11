package db

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

type DBLayer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Ping(ctx context.Context) error
	Close()
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Client struct {
	db DBLayer
}

func NewClient(db DBLayer) *Client {
	return &Client{db: db}
}

func (p *Client) GetDB() DBLayer {
	return p.db
}

func (p *Client) ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}
	return p.db.Exec(ctx, q.QueryRaw, args...)
}

func (p *Client) QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}
	return p.db.Query(ctx, q.QueryRaw, args...)
}

func (p *Client) QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.db.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *Client) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func (p *Client) Close() {
	p.db.Close()
}

func (p *Client) ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, row)
}

func (p *Client) ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, row)
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
