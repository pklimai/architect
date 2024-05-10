package architect

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	gw_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func (a *App) initGRPCServer(services ...service) {
	if len(services) == 0 {
		return
	}

	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(a.options.UnaryInterseptors...),
	)

	for _, s := range services {
		s.RegisterGRPC(a.grpcServer)
	}

	reflection.Register(a.grpcServer)
}

func (a *App) initHTTPServer(services ...service) error {
	a.httpServer = chi.NewRouter()

	a.httpServer.Use(cors.Handler(a.options.corsOptions))

	gatewayMux := gw_runtime.NewServeMux()

	for _, s := range services {
		// TODO: Need to use RegisterGateway, but middlewares are problems
		// coz gateway use gRPC server directly without them.
		if err := s.RegisterGatewayEndpoint(
			context.Background(),
			gatewayMux,
			fmt.Sprintf(formatAddress, a.options.host, a.options.portGRPC),
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			},
		); err != nil {
			return fmt.Errorf("RegisterGatewayEndpoint: %w", err)
		}
	}

	a.httpServer.Handle("/*", gatewayMux)

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

	a.swaggerServer = http.FileServer(http.FS(subFS))

	return nil
}
