package env

import (
	"errors"
	"net"
	"os"
)

const (
	GRPCHost = "GRPC_HOST"
	GRPCPort = "GRPC_PORT"
)

type GRPCConfig struct {
	addr string
}

func NewGRPCConfig() (*GRPCConfig, error) {
	host := os.Getenv(GRPCHost)
	if host == "" {
		return nil, errors.New("host is empty")
	}

	port := os.Getenv(GRPCPort)
	if port == "" {
		return nil, errors.New("port is empty")
	}

	return &GRPCConfig{addr: net.JoinHostPort(host, port)}, nil
}

func (g *GRPCConfig) Addr() string {
	return g.addr
}
