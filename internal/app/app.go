package app

import (
	"context"
	"flag"
	"log/slog"
	"net"

	"github.com/Str1m/auth/internal/config"
	"github.com/Str1m/auth/internal/lib/logger/sl"
	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *ServiceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	defer func() {
		a.serviceProvider.cls.CloseAll()
		a.serviceProvider.cls.Wait()
	}()
	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPC,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	var cfgPath string
	flag.StringVar(&cfgPath, "config-path", ".env", "path to config file")
	flag.Parse()

	err := config.MustLoad(cfgPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPC(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.UserAPIImpl(ctx))
	a.serviceProvider.GetCloser().Add(func() error {
		a.grpcServer.GracefulStop()
		return nil
	})
	return nil
}

func (a *App) runGRPCServer() error {
	a.serviceProvider.GetLog().Info("server listening", slog.String("Addr", a.serviceProvider.GRPCConfig().Address()))

	l, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		a.serviceProvider.GetLog().Error("failed to listen", sl.Err(err))
		return err
	}

	err = a.grpcServer.Serve(l)
	if err != nil {
		return err
	}
	return nil
}
