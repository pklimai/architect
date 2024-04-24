package architect

import (
	"embed"
	"time"
)

const (
	gracefulDelay   = 5 * time.Second
	gracefulTimeout = 10 * time.Second

	readHeaderTimeout = 5 * time.Second
)

type options struct {
	host        string
	portHTTP    uint
	portSwagger uint
	portGRPC    uint

	readHeaderTimeout time.Duration

	gracefulDelay   time.Duration
	gracefulTimeout time.Duration

	swaggerFS embed.FS // TODO: move swagger to architect
}

func initOptions(appSettings AppSettings) options {
	opts := options{
		host:              appSettings.Host,
		portHTTP:          appSettings.PortHTTP,
		portSwagger:       appSettings.PortSwagger,
		portGRPC:          appSettings.PortGRPC,
		readHeaderTimeout: readHeaderTimeout,
		gracefulDelay:     gracefulDelay,
		gracefulTimeout:   gracefulTimeout,
		swaggerFS:         appSettings.SwaggerFS,
	}

	if opts.host == "" {
		opts.host = "localhost"
	}

	if opts.portHTTP == 0 {
		opts.portHTTP = 7000
	}

	if opts.portSwagger == 0 {
		opts.portSwagger = 7001
	}

	if opts.portGRPC == 0 {
		opts.portGRPC = 7002
	}

	return opts
}
