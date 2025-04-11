package postgres

import (
	"context"
	"github.com/Str1m/auth/internal/client/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

//type DB interface {
//	SQLExecer
//	Pinger
//	Close()
//}
//
//type SQLExecer interface {
//	NamedExecer
//	QueryExecer
//}
//
//type NamedExecer interface {
//	ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error
//	ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error
//}
//
//type QueryExecer interface {
//	ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error)
//	QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error)
//	QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row
//}
//
//type Pinger interface {
//	Ping(ctx context.Context) error
//}
//
//type Client interface {
//	DB() DB
//	Close() error
//}

type pgClient struct {
	masterDB db.DB
}

func NewPGClient(ctx context.Context, dsn string) (db.Client, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &pgClient{
		masterDB: &pg{db: db},
	}, nil
}

func (p *pgClient) DB() db.DB {
	return p.masterDB
}

func (p *pgClient) Close() error {
	if p.masterDB != nil {
		p.masterDB.Close()
	}
	return nil
}
