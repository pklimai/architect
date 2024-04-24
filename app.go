package architect

import (
	"fmt"
	"net/http"
	"syscall"

	"gitlab.com/zigal0/architect/pkg/closer"
	"gitlab.com/zigal0/architect/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	options options

	closer *closer.Closer

	listeners *listeners

	grpcServer    *grpc.Server
	httpServer    *http.Server
	swaggerServer *http.Server
}

func NewApp(appSettings AppSettings) (*App, error) {
	logLevel, err := logger.ParseLevel(appSettings.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(logLevel)

	options := initOptions(appSettings)

	listeners, err := newListeners(options)
	if err != nil {
		return nil, fmt.Errorf("newListeners: %w", err)
	}

	return &App{
		options:   options,
		listeners: listeners,
		closer:    closer.New(syscall.SIGTERM, syscall.SIGINT),
	}, nil
}
