package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"syscall"

	userAPI "github.com/Str1m/auth/internal/api/user"
	"github.com/Str1m/auth/internal/closer"
	"github.com/Str1m/auth/internal/config/env"
	"github.com/Str1m/auth/internal/lib/logger/handlers/slogpretty"
	modelService "github.com/Str1m/auth/internal/model"
	"github.com/Str1m/auth/internal/service/user"
	"github.com/Str1m/auth/internal/storage/users/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Storage interface {
	Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

type Service interface {
	Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name *string, email *string) error
	Delete(ctx context.Context, id int64) error
}

type StorageConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}

type ServiceProvider struct {
	cls *closer.Closer
	log *slog.Logger

	repoConfig StorageConfig
	grpcConfig GRPCConfig

	pgPool         *pgxpool.Pool
	userRepository Storage

	userService Service
	userAPI     *userAPI.Implementation
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) Closer() *closer.Closer {
	if s.cls == nil {
		s.cls = closer.New(os.Interrupt, syscall.SIGTERM)
	}

	return s.cls
}

func (s *ServiceProvider) Log() *slog.Logger {
	if s.log == nil {
		env := os.Getenv("ENV")
		switch env {
		case envLocal:
			s.log = setupPrettySlog()
		case envDev:
			s.log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envProd:
			s.log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
		}
	}
	return s.log
}

func (s *ServiceProvider) RepoConfig() StorageConfig {
	if s.repoConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.repoConfig = cfg
	}
	return s.repoConfig
}

func (s *ServiceProvider) GRPCConfig() GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *ServiceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.RepoConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}
		if err = pool.Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		s.Closer().Add(func() error {
			pool.Close()
			return nil
		})
		s.pgPool = pool
	}
	return s.pgPool
}

func (s *ServiceProvider) UserRepository(ctx context.Context) Storage {
	if s.userRepository == nil {
		s.userRepository = postgres.NewRepository(s.PGPool(ctx))
	}

	return s.userRepository
}

func (s *ServiceProvider) UserService(ctx context.Context) Service {
	if s.userService == nil {
		s.userService = user.NewService(s.Log(), s.UserRepository(ctx))
	}
	return s.userService
}

func (s *ServiceProvider) UserAPIImpl(ctx context.Context) *userAPI.Implementation {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewImplementation(s.UserService(ctx))
	}

	return s.userAPI
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
