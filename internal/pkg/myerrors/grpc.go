package myerrors

import "google.golang.org/grpc/status"

type GrpcErr struct{}

func (g *GrpcErr) New(err error) error {
	unwrappedErr := UnWrapErr(err)

	code := unwrappedErr.GetGrpcCode()

	return status.New(code, unwrappedErr.message).Err()
}
