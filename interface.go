package architect

import (
	"context"

	gw_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type service interface {
	RegisterGRPC(server *grpc.Server)
	RegisterGatewayEndpoint(
		ctx context.Context,
		mux *gw_runtime.ServeMux,
		endpoint string,
		dialOptions []grpc.DialOption,
	) error
}
