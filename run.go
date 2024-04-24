package architect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/zigal0/architect/pkg/closer"
	"gitlab.com/zigal0/architect/pkg/logger"
)

func (a *App) Run(services ...service) error {
	// need to init first due to hack in initHTTPServer.
	a.initGRPCServer(services...)
	a.runGRPC()

	if err := a.initHTTPServer(services...); err != nil {
		return fmt.Errorf("initHTTPServer: %w", err)
	}

	a.runHTTP()

	if err := a.initSwaggerServer(); err != nil {
		return fmt.Errorf("initSwaggerServer: %w", err)
	}

	a.runSwagger()

	logger.Errorf(
		"App started on host = %s, ports: HTTP = %d, Swagger = %d, GRPC = %d",
		a.options.host,
		a.options.portHTTP,
		a.options.portSwagger,
		a.options.portGRPC,
	)

	a.closer.Wait()

	closer.CloseAll()

	return nil
}

func (a *App) runGRPC() {
	if a.grpcServer == nil {
		return
	}

	go func() {
		if err := a.grpcServer.Serve(a.listeners.grpc); err != nil {
			logger.Errorf("grpcServer.Serve: %v", err)

			a.closer.CloseAll()
		}
	}()

	a.closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), a.options.gracefulTimeout)
		defer cancel()

		logger.Warn("grpc server is waiting for traffic stop")
		time.Sleep(a.options.gracefulDelay)
		logger.Warn("grpc server shutting down")

		done := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()

			close(done)
		}()

		select {
		case <-done:
			logger.Warn("grpc server was gracefully stopped")
		case <-ctx.Done():
			a.grpcServer.Stop()

			return fmt.Errorf("grpc server was forcibly stopped: %w", ctx.Err())
		}

		return nil
	})
}

func (a *App) runHTTP() {
	if a.httpServer == nil {
		return
	}

	go func() {
		if err := a.httpServer.Serve(a.listeners.http); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("httpServer.Serve: %v", err)

			a.closer.CloseAll()
		}
	}()

	a.closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), a.options.gracefulTimeout)
		defer cancel()

		logger.Warn("http server is waiting for traffic stop")
		time.Sleep(a.options.gracefulDelay)
		logger.Warn("http server is shutting down")

		a.httpServer.SetKeepAlivesEnabled(false)

		if err := a.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("http server failed to shutdown: %w", err)
		}

		logger.Warn("http server was gracefully stopped")

		return nil
	})
}

func (a *App) runSwagger() {
	if a.swaggerServer == nil {
		return
	}

	go func() {
		if err := a.swaggerServer.Serve(a.listeners.swagger); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("swaggerServer.Serve: %v", err)

			a.closer.CloseAll()
		}
	}()

	a.closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), a.options.gracefulTimeout)
		defer cancel()

		logger.Warn("swagger server is waiting for traffic stop")
		time.Sleep(a.options.gracefulDelay)
		logger.Warn("swagger server is shutting down")

		a.swaggerServer.SetKeepAlivesEnabled(false)

		if err := a.swaggerServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("swagger server failed to shutdown: %w", err)
		}

		logger.Warn("swagger server was gracefully stopped")

		return nil
	})
}
