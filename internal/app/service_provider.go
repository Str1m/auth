package app

import (
	"context"
	modelService "github.com/Str1m/auth/internal/model"
	"github.com/Str1m/auth/internal/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"os"
	"syscall"

	userAPI "github.com/Str1m/auth/internal/api/user"
	"github.com/Str1m/auth/internal/closer"
	"github.com/Str1m/auth/internal/config/env"
	"github.com/Str1m/auth/internal/lib/logger/handlers/slogpretty"
	"github.com/Str1m/auth/internal/storage/users/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type StorageConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}

type Service interface {
	Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name *string, email *string) error
	Delete(ctx context.Context, id int64) error
}

type Storage interface {
	Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

type ServiceProvider struct {
	cls *closer.Closer
	log *slog.Logger

	storageConfig StorageConfig
	grpcConfig    GRPCConfig

	dbClient *postgres.ClientPG
	dbLayer  Storage

	userService Service
	userAPI     *userAPI.Implementation
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) GetCloser() *closer.Closer {
	if s.cls == nil {
		s.cls = closer.New(os.Interrupt, syscall.SIGTERM)
	}

	return s.cls
}

func (s *ServiceProvider) GetLog() *slog.Logger {
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

func (s *ServiceProvider) GetStorageConfig() StorageConfig {
	if s.storageConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.storageConfig = cfg
	}
	return s.storageConfig
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

func (s *ServiceProvider) GetDBClient(ctx context.Context) *postgres.ClientPG {
	if s.dbClient == nil {
		p, err := pgxpool.New(ctx, s.GetStorageConfig().DSN())
		if err != nil {
			log.Fatalln("err")
		}

		s.dbClient = postgres.NewClientPG(p)
	}

	return s.dbClient
}

func (s *ServiceProvider) GetDBLayer(ctx context.Context) Storage {
	if s.dbLayer == nil {
		s.dbLayer = postgres.NewStoragePG(s.GetDBClient(ctx))
	}
	return s.dbLayer
}

func (s *ServiceProvider) UserService(ctx context.Context) Service {
	if s.userService == nil {
		s.userService = user.NewService(s.GetLog(), s.GetDBLayer(ctx))
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
