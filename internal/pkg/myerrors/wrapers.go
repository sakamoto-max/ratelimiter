package myerrors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type WrappedErr struct {
	message string
	errType string
}

func (e *WrappedErr) Error() string {
	return e.message
}

func (e *WrappedErr) GetGrpcCode() codes.Code {
	switch e.errType {
	case InternalServerErr:
		return codes.Internal
	case BadRequestErr:
		return codes.InvalidArgument
	case NotFoundErr:
		return codes.NotFound
	case AlreadyExistsErr:
		return codes.AlreadyExists
	case UnauthorizedErr:
		return codes.Unauthenticated
	}

	return codes.Unknown
}
func (e *WrappedErr) GetHttpCode() int {
	switch e.errType {
	case InternalServerErr:
		return http.StatusInternalServerError
	case BadRequestErr:
		return http.StatusBadRequest
	case NotFoundErr:
		return http.StatusNotFound
	case AlreadyExistsErr:
		return http.StatusConflict
	case UnauthorizedErr:
		return http.StatusUnauthorized
	}

	return 0
}

func WrapErr(err error, errtype string) error {
	return &WrappedErr{
		message: err.Error(),
		errType: errtype,
	}
}

func UnWrapErr(err error) *WrappedErr {
	unwrappedErr, ok := err.(*WrappedErr)
	if !ok {
		return nil
	}

	return unwrappedErr

}

var (
	InternalServerErr string = "internalServer"
	BadRequestErr     string = "badRequest"
	NotFoundErr       string = "notFound"
	AlreadyExistsErr  string = "alreadyExists"
	UnauthorizedErr   string = "unauthorized"
)
