package business_error

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/zigal0/architect/pkg/logger"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor
// for handling business errors and logging errors.
func UnaryServerInterceptor(enableLog bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return resp, err
		}

		err = fmt.Errorf(constructServiceAndHandlerFromat(info.FullMethod), err)

		grpcErr := ToGRPCError(err)

		if enableLog {
			logError(grpcErr.Code(), err)
		}

		return resp, grpcErr.Err()
	}
}

func constructServiceAndHandlerFromat(rawMethod string) (res string) {
	defer func() {
		res += ": %w"
	}()

	serviceMethodFromat := "UnrecognizedService.UnrecognizedMethod"
	fullMethodParts := strings.Split(rawMethod, ".")

	if len(fullMethodParts) == 0 {
		return serviceMethodFromat
	}

	serviceMethodFromat = fullMethodParts[len(fullMethodParts)-1]

	return strings.Replace(serviceMethodFromat, "/", ".", 1)
}

func logError(code codes.Code, err error) {
	switch getLogLevelByStatusCode(code) {
	case zapcore.DebugLevel:
		logger.Debug(err)
	case zapcore.InfoLevel:
		logger.Info(err)
	case zapcore.WarnLevel:
		logger.Warn(err)
	default:
		logger.Error(err)
	}
}
