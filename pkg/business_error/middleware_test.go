package business_error_test

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/zigal0/architect/pkg/business_error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_HandleErrorMiddleware(t *testing.T) {
	t.Parallel()

	var (
		handlerWrap = "TestService/TestHandler"

		info = &grpc.UnaryServerInfo{
			FullMethod: "some.internal.grpc.info" + handlerWrap,
		}
	)

	t.Run("no FullMethod", func(t *testing.T) {
		t.Parallel()

		// arrange
		f := setUp(t)

		var (
			errFromHandler = business_error.New(errForTest, errMsgForTest, business_error.Canceled)

			handler = grpc.UnaryHandler(
				func(context.Context, interface{}) (interface{}, error) {
					return nil, errFromHandler
				},
			)
		)

		mw := business_error.HandleErrorMiddleware(true)

		// act
		res, err := mw(f.ctx, nil, &grpc.UnaryServerInfo{}, handler)

		// assert
		f.Empty(res, "response")
		f.Equal(codes.Canceled, status.Code(err), "error code")
	})

	t.Run("exist error", func(t *testing.T) {
		t.Parallel()

		// arrange
		f := setUp(t)

		var (
			errFromHandler = business_error.New(errForTest, errMsgForTest, business_error.DataLoss)

			handler = grpc.UnaryHandler(
				func(context.Context, interface{}) (interface{}, error) {
					return nil, errFromHandler
				},
			)
		)

		mw := business_error.HandleErrorMiddleware(true)

		// act
		res, err := mw(f.ctx, nil, info, handler)

		// assert
		f.Empty(res, "response")
		f.Equal(codes.DataLoss, status.Code(err), "error code")
	})

	t.Run("no business wrap", func(t *testing.T) {
		t.Parallel()

		// arrange
		f := setUp(t)

		var (
			errFromHandler = errForTest

			handler = grpc.UnaryHandler(
				func(context.Context, interface{}) (interface{}, error) {
					return nil, errFromHandler
				},
			)
		)

		mw := business_error.HandleErrorMiddleware(true)

		// act
		res, err := mw(f.ctx, nil, info, handler)

		// assert
		f.Empty(res, "response")
		f.Equal(codes.Internal, status.Code(err), "error code")
	})

	t.Run("additional wrap", func(t *testing.T) {
		t.Parallel()

		// arrange
		f := setUp(t)

		var (
			errFromHandler = fmt.Errorf(
				"wrap: %w",
				business_error.New(errForTest, errMsgForTest, business_error.AlreadyExists),
			)

			handler = grpc.UnaryHandler(
				func(context.Context, interface{}) (interface{}, error) {
					return nil, errFromHandler
				},
			)
		)

		mw := business_error.HandleErrorMiddleware(true)

		// act
		res, err := mw(f.ctx, nil, info, handler)

		// assert
		f.Empty(res, "response")
		f.Equal(codes.AlreadyExists, status.Code(err), "error code")
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		// arrange
		f := setUp(t)

		var (
			handler = grpc.UnaryHandler(
				func(context.Context, interface{}) (interface{}, error) {
					return nil, nil
				},
			)
		)

		mw := business_error.HandleErrorMiddleware(true)

		// act
		res, err := mw(f.ctx, nil, info, handler)

		// assert
		f.Empty(res, "response")
		f.NoError(err, "error")
	})
}
