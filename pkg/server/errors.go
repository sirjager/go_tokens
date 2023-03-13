package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unAuthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
