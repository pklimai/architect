package architect

import (
	"embed"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/cors"
	"google.golang.org/grpc"
)

const (
	hostDefault = "127.0.0.1"

	portHTTPDefault    = 7000
	portSwaggerDefault = 7001
	portGRPCDefault    = 7002

	gracefulDelayDefault   = 5 * time.Second
	gracefulTimeoutDefault = 10 * time.Second

	readHeaderTimeoutDefault = 5 * time.Second

	defaultPreflightMaxAge = 3600
)

var (
	defaultAllowedMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,

		// For CORS preflight requests
		http.MethodOptions,
	}

	allHTTPMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
)

type Options struct {
	corsOptions cors.Options

	host        string
	portHTTP    uint
	portSwagger uint
	portGRPC    uint

	readHeaderTimeout time.Duration

	gracefulDelay   time.Duration
	gracefulTimeout time.Duration

	swaggerFS embed.FS // TODO: move swagger to architect

	UnaryInterseptors []grpc.UnaryServerInterceptor
}

func initOptions(appSettings AppSettings, optionAppliers ...OptionApplier) (*Options, error) {
	defaultCORSOptions := initDefaultCORSOptions()

	opts := &Options{
		corsOptions:       defaultCORSOptions,
		host:              appSettings.Host,
		portHTTP:          appSettings.PortHTTP,
		portSwagger:       appSettings.PortSwagger,
		portGRPC:          appSettings.PortGRPC,
		readHeaderTimeout: readHeaderTimeoutDefault,
		gracefulDelay:     gracefulDelayDefault,
		gracefulTimeout:   gracefulTimeoutDefault,
		swaggerFS:         appSettings.SwaggerFS,
	}

	if opts.host == "" {
		opts.host = hostDefault
	}

	if opts.portHTTP == 0 {
		opts.portHTTP = portHTTPDefault
	}

	if opts.portSwagger == 0 {
		opts.portSwagger = portSwaggerDefault
	}

	if opts.portGRPC == 0 {
		opts.portGRPC = portGRPCDefault
	}

	for _, optApplier := range optionAppliers {
		err := optApplier.Apply(opts)
		if err != nil {
			return nil, fmt.Errorf("optApplier.Apply: %w", err)
		}
	}

	return opts, nil
}

func initDefaultCORSOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowOriginFunc:    nil,
		AllowedMethods:     defaultAllowedMethods,
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{},
		AllowCredentials:   false,
		MaxAge:             defaultPreflightMaxAge,
		OptionsPassthrough: false,
		Debug:              false,
	}
}

// TODO: move all OptionApplier to one file.
type OptionApplier interface {
	Apply(*Options) error
}

type optionApplier func(*Options) error

func (oa optionApplier) Apply(opts *Options) error {
	return oa(opts)
}

func WithUnaryInterseptor(iterceptor grpc.UnaryServerInterceptor) OptionApplier {
	return optionApplier(func(opts *Options) error {
		if iterceptor != nil {
			opts.UnaryInterseptors = append(opts.UnaryInterseptors, iterceptor)
		}

		return nil
	})
}

func WithCORSAllowedOrigins(origins []string) OptionApplier {
	return optionApplier(func(opts *Options) error {
		opts.corsOptions.AllowedOrigins = origins

		return nil
	})
}

func WithCORSAllowedMethods(methods []string) OptionApplier {
	return optionApplier(func(opts *Options) error {
		if err := validateHTTPMethods(methods); err != nil {
			return err
		}

		opts.corsOptions.AllowedMethods = methods

		return nil
	})
}

func validateHTTPMethods(methods []string) error {
	for _, method := range methods {
		if !slices.Contains(allHTTPMethods, method) {
			return fmt.Errorf("invslid name of http method: %s", method)
		}
	}

	return nil
}

func WithCORSAllowedHeaders(headers []string) OptionApplier {
	return optionApplier(func(opts *Options) error {
		opts.corsOptions.AllowedHeaders = headers

		return nil
	})
}

func WithCORSExposedHeaders(headers []string) OptionApplier {
	return optionApplier(func(opts *Options) error {
		opts.corsOptions.ExposedHeaders = headers

		return nil
	})
}
