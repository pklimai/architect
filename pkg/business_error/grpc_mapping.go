package business_error

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError converts Error into gRPC error.
func ToGRPCError(err error) *status.Status {
	if err == nil {
		return nil
	}

	if businessError := new(Error); errors.As(err, &businessError) {
		return status.New(
			ToGRPCCode(businessError.GetCode()),
			businessError.GetMessage(),
		)
	}

	return status.New(codes.Internal, "Internal service error")
}

// ToGRPCCode converts Error Code into gRPC Code.
func ToGRPCCode(code Code) codes.Code {
	var grpcCode codes.Code

	switch code {
	case OK:
		grpcCode = codes.OK
	case Canceled:
		grpcCode = codes.Canceled
	case Unknown:
		grpcCode = codes.Unknown
	case InvalidArgument:
		grpcCode = codes.InvalidArgument
	case DeadlineExceeded:
		grpcCode = codes.DeadlineExceeded
	case NotFound:
		grpcCode = codes.NotFound
	case AlreadyExists:
		grpcCode = codes.AlreadyExists
	case PermissionDenied:
		grpcCode = codes.PermissionDenied
	case ResourceExhausted:
		grpcCode = codes.ResourceExhausted
	case FailedPrecondition:
		grpcCode = codes.FailedPrecondition
	case Aborted:
		grpcCode = codes.Aborted
	case OutOfRange:
		grpcCode = codes.OutOfRange
	case Unimplemented:
		grpcCode = codes.Unimplemented
	case Internal:
		grpcCode = codes.Internal
	case Unavailable:
		grpcCode = codes.Unavailable
	case DataLoss:
		grpcCode = codes.DataLoss
	case Unauthenticated:
		grpcCode = codes.Unauthenticated
	default:
		grpcCode = codes.Unknown
	}

	return grpcCode
}

// FromGRPCCode converts gRPC Code into Error Code.
func FromGRPCCode(code codes.Code) Code {
	var businessCode Code

	switch code {
	case codes.OK:
		businessCode = OK
	case codes.Canceled:
		businessCode = Canceled
	case codes.Unknown:
		businessCode = Unknown
	case codes.InvalidArgument:
		businessCode = InvalidArgument
	case codes.DeadlineExceeded:
		businessCode = DeadlineExceeded
	case codes.NotFound:
		businessCode = NotFound
	case codes.AlreadyExists:
		businessCode = AlreadyExists
	case codes.PermissionDenied:
		businessCode = PermissionDenied
	case codes.ResourceExhausted:
		businessCode = ResourceExhausted
	case codes.FailedPrecondition:
		businessCode = FailedPrecondition
	case codes.Aborted:
		businessCode = Aborted
	case codes.OutOfRange:
		businessCode = OutOfRange
	case codes.Unimplemented:
		businessCode = Unimplemented
	case codes.Internal:
		businessCode = Internal
	case codes.Unavailable:
		businessCode = Unavailable
	case codes.DataLoss:
		businessCode = DataLoss
	case codes.Unauthenticated:
		businessCode = Unauthenticated
	default:
		businessCode = Unknown
	}

	return businessCode
}
