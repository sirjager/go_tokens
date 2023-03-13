package server

import (
	"context"

	rpc "github.com/sirjager/rpcs/tokens/go"
)

func (srvr *RPCTokensServer) TokensWelcome(ctx context.Context, req *rpc.TokensWelcomeRequest) (*rpc.TokensWelcomeResponse, error) {
	res := &rpc.TokensWelcomeResponse{Message: "Welcome to Tokens Api"}
	return res, nil
}
