package architect

import (
	"fmt"
	"net"
)

const (
	formatAddress = "%s:%d"
)

type listeners struct {
	http    net.Listener
	grpc    net.Listener
	swagger net.Listener
}

const (
	formatErrFailedToListen = "failed to init %s listener for port %d: %w"
)

func newListeners(options *Options) (*listeners, error) {
	httpListener, err := net.Listen("tcp", fmt.Sprintf(
		formatAddress, options.host, options.portHTTP,
	))
	if err != nil {
		return nil, fmt.Errorf(formatErrFailedToListen, "http", options.portHTTP, err)
	}

	swaggerListener, err := net.Listen("tcp", fmt.Sprintf(
		formatAddress, options.host, options.portSwagger,
	))
	if err != nil {
		return nil, fmt.Errorf(formatErrFailedToListen, "swagger", options.portHTTP, err)
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(
		formatAddress, options.host, options.portGRPC,
	))
	if err != nil {
		return nil, fmt.Errorf(formatErrFailedToListen, "grpc", options.portHTTP, err)
	}

	return &listeners{
		http:    httpListener,
		grpc:    grpcListener,
		swagger: swaggerListener,
	}, nil
}
