package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func GRPCMustHave(header string, expectedValue string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		values, err := GetHeaders(ctx, header)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if len(values) != 1 {
			return nil, status.Errorf(codes.Internal, "invalid %s ", header)
		}
		if expectedValue == values[0] {
			res, err = handler(ctx, req)
			return res, err
		}

		return nil, status.Errorf(codes.Internal, "access denined: invalid %s", header)
	}
}

func HTTPMustHave(handler http.Handler, header string, expectedValue string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		value := req.Header.Get(header)
		if value == "" {
			errResponse := ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    fmt.Sprintf("missing %s header", header),
			}

			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(errResponse.StatusCode)

			err := json.NewEncoder(res).Encode(errResponse)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		if value != expectedValue {
			errResponse := ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    fmt.Sprintf("invalid %s header", header),
			}
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(errResponse.StatusCode)

			err := json.NewEncoder(res).Encode(errResponse)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		handler.ServeHTTP(res, req)
	})
}
