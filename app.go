package architect

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/go-chi/chi/v5"
	"gitlab.com/zigal0/architect/pkg/closer"
	"gitlab.com/zigal0/architect/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	options *Options

	closer *closer.Closer

	listeners *listeners

	httpServer    chi.Router
	swaggerServer http.Handler
	grpcServer    *grpc.Server
}

func NewApp(appSettings AppSettings, optionAppliers ...OptionApplier) (*App, error) {
	logLevel, err := logger.ParseLevel(appSettings.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(logLevel)

	options, err := initOptions(appSettings, optionAppliers...)
	if err != nil {
		return nil, fmt.Errorf("initOptions: %w", err)
	}

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
