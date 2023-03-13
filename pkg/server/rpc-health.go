package server

import (
	"context"
	"time"

	rpc "github.com/sirjager/rpcs/tokens/go"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srvr *RPCTokensServer) TokensHealth(ctx context.Context, req *rpc.TokensHealthRequest) (*rpc.TokensHealthResponse, error) {
	res := &rpc.TokensHealthResponse{
		Status:    "UP",
		Timestamp: timestamppb.New(time.Now()),
		Protected: srvr.config.ApiSecret != "",
		Uptime:    durationpb.New(time.Since(srvr.startTime)),
		Started:   timestamppb.New(srvr.startTime),
	}
	return res, nil
}
