package app

import (
	"context"
	"github.com/Str1m/auth/internal/client/db"
	dbPG "github.com/Str1m/auth/internal/client/db/postgres"
	"github.com/Str1m/auth/internal/client/db/transaction"
	"github.com/Str1m/auth/internal/service"
	"github.com/Str1m/auth/internal/storage"
	"log"
	"log/slog"
	"os"
	"syscall"

	userAPI "github.com/Str1m/auth/internal/api/user"
	"github.com/Str1m/auth/internal/closer"
	"github.com/Str1m/auth/internal/config/env"
	"github.com/Str1m/auth/internal/lib/logger/handlers/slogpretty"
	"github.com/Str1m/auth/internal/service/user"
	"github.com/Str1m/auth/internal/storage/users/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

//	type Storage interface {
//		Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
//		Get(ctx context.Context, id int64) (*modelService.User, error)
//		Update(ctx context.Context, id int64, name, email *string) error
//		Delete(ctx context.Context, id int64) error
//	}
//
//	type Service interface {
//		Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error)
//		Get(ctx context.Context, id int64) (*modelService.User, error)
//		Update(ctx context.Context, id int64, name *string, email *string) error
//		Delete(ctx context.Context, id int64) error
//	}
type StorageConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}

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

type ServiceProvider struct {
	cls *closer.Closer
	log *slog.Logger

	repoConfig StorageConfig
	grpcConfig GRPCConfig

	txManager      db.TxManager
	dbClient       db.Client
	userRepository storage.Repository

	userService service.Service
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

func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *ServiceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := dbPG.NewPGClient(ctx, s.RepoConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}
		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		s.Closer().Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}

func (s *ServiceProvider) UserRepository(ctx context.Context) storage.Repository {
	if s.userRepository == nil {
		s.userRepository = postgres.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *ServiceProvider) UserService(ctx context.Context) service.Service {
	if s.userService == nil {
		s.userService = user.NewService(s.Log(), s.UserRepository(ctx), s.TxManager(ctx))
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
