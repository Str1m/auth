package main

import (
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/Str1m/auth/internal/config"
	"github.com/Str1m/auth/internal/config/env"
	"github.com/Str1m/auth/internal/grpc/auth"
	"github.com/Str1m/auth/internal/lib/logger/handlers/slogpretty"
	"github.com/Str1m/auth/internal/lib/logger/sl"
	"github.com/Str1m/auth/internal/storage/postgres"
	desc "github.com/Str1m/auth/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config-path", ".env", "path to config file")
	flag.Parse()

	config.MustLoad(cfgPath)

	log := setupLogger()

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Error("failed to get grpc config", sl.Err(err))
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Error("failed to get postgres config", sl.Err(err))
	}

	_, err = postgres.New(pgConfig.DSN())
	if err != nil {
		log.Error("failed to connect to db", sl.Err(err))
	}

	l, err := net.Listen("tcp", grpcConfig.Addr())
	if err != nil {
		log.Error("failed to listen", sl.Err(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterAuthV1Server(s, &auth.Server{})

	log.Info("server listening", slog.String("Addr", grpcConfig.Addr()))

	if err = s.Serve(l); err != nil {
		log.Error("failed to serve", sl.Err(err))
	}
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	env := os.Getenv("ENV")
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
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
