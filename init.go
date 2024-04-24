package architect

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	gw_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.com/zigal0/architect/pkg/business_error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func (a *App) initGRPCServer(services ...service) {
	if len(services) == 0 {
		return
	}

	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_validator.UnaryServerInterceptor(),
			business_error.HandleErrorMiddleware(true),
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	for _, s := range services {
		s.RegisterGRPC(a.grpcServer)
	}

	reflection.Register(a.grpcServer)
}

func (a *App) initHTTPServer(services ...service) error {
	gatewayMux := gw_runtime.NewServeMux()

	for _, s := range services {
		// Need to use RegisterGateway, but middlewares are problems.
		if err := s.RegisterGatewayEndpoint(
			context.Background(),
			gatewayMux,
			fmt.Sprintf(formatAddress, a.options.host, a.options.portGRPC),
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			},
		); err != nil {
			return fmt.Errorf("RegisterGateway: %w", err)
		}
	}

	a.httpServer = &http.Server{
		ReadHeaderTimeout: a.options.readHeaderTimeout,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			gatewayMux.ServeHTTP(w, r)
		}),
	}

	return nil
}

// TODO: dirty code, need refactoring
func (a *App) initSwaggerServer() error {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return fmt.Errorf("mime.AddExtensionType: %w", err)
	}

	// Use subdirectory in embedded files
	subFS, err := fs.Sub(a.options.swaggerFS, "src")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}

	handler := http.FileServer(http.FS(subFS))

	a.swaggerServer = &http.Server{
		ReadHeaderTimeout: a.options.readHeaderTimeout,
		Handler:           handler,
	}

	return nil
}
