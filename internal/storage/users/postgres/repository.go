package postgres

import (
	"context"
	"fmt"

	modelService "github.com/Str1m/auth/internal/model"

	"github.com/Str1m/auth/internal/storage"
	"github.com/Str1m/auth/internal/storage/users/converter"
	modelRepo "github.com/Str1m/auth/internal/storage/users/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName = "users"

	idColumn           = "id"
	nameColumn         = "name"
	emailColumn        = "email"
	passwordHashColumn = "password_hash"
	roleColumn         = "role"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Close() {
	r.db.Close()
}

func (r *Repo) Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error) {
	const op = "repository.users.Create"

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordHashColumn, roleColumn).
		Values(info.Name, info.Email, hashedPassword, info.Role).Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int64) (*modelService.User, error) {
	const op = "repository.users.Get"

	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	var user modelRepo.User
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToUserFromStorage(&user), nil
}

func (r *Repo) Update(ctx context.Context, id int64, name, email *string) error {
	const op = "repository.users.Update"

	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	if name != nil {
		builder = builder.Set(nameColumn, *name)
	}

	if email != nil {
		builder = builder.Set(emailColumn, *email)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	const op = "repository.users.Delete"

	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}
