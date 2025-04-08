package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Str1m/auth/internal/models"
	desc "github.com/Str1m/auth/pkg/auth_v1"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	_ = s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, name string, email string,
	passHash []byte, role desc.Role) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, name, email, passHash, int32(role)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id int64, name, email *string) error {
	const op = "storage.postgres.UpdateUser"

	if name == nil && email == nil {
		return nil
	}

	query := `UPDATE users SET `
	args := []interface{}{}
	argIndex := 1

	if name != nil {
		query += fmt.Sprintf("name = $%d, ", argIndex)
		args = append(args, *name)
		argIndex++
	}
	if email != nil {
		query += fmt.Sprintf("email = $%d, ", argIndex)
		args = append(args, *email)
		argIndex++
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, id)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) User(ctx context.Context, email string) (models.UserInfo, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u models.UserInfo
	var role int32
	var createdAt, updatedAt time.Time

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &role, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("%s: user not found: %w", op, err)
		}
		return u, fmt.Errorf("%s: %w", op, err)
	}

	u.Role = desc.Role(role)
	u.CreatedAt = timestamppb.New(createdAt)
	u.UpdatedAt = timestamppb.New(updatedAt)

	return u, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int64) error {
	const op = "storage.postgres.DeleteUser"

	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no user with id %d", op, id)
	}

	return nil
}
