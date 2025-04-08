package env

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	pgUser     = "POSTGRES_USER"
	pgPassword = "POSTGRES_PASSWORD"
	pgHost     = "POSTGRES_HOST"
	pgPort     = "POSTGRES_PORT"
	pgDBName   = "POSTGRES_DB"
)

type PostgresConfig struct {
	dsn string
}

func NewPGConfig() (*PostgresConfig, error) {
	user := os.Getenv(pgUser)
	password := os.Getenv(pgPassword)
	host := os.Getenv(pgHost)
	port := os.Getenv(pgPort)
	dbname := os.Getenv(pgDBName)

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		return nil, errors.New("some required environment variables are missing")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		user, password, net.JoinHostPort(host, port), dbname)

	return &PostgresConfig{dsn: dsn}, nil
}

func (p *PostgresConfig) DSN() string {
	return p.dsn
}
