package server

import (
	"context"
	"strings"

	"github.com/sirjager/go_tokens/pkg/tokens"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpc "github.com/sirjager/rpcs/tokens/go"
)

func (srvr *RPCTokensServer) TokensVerify(ctx context.Context, req *rpc.TokensVerifyRequest) (*rpc.TokensVerifyResponse, error) {
	// This only validates if the request is authorized or not
	if err := srvr.authorize(ctx); err != nil {
		return nil, unAuthenticatedError(err)
	}

	var builder tokens.TokenBuilder
	switch strings.ToLower(srvr.config.TokenBuilder) {
	case tokens.Paseto:
		builder = srvr.pasetoBuilder
	case tokens.Jwt:
		builder = srvr.jwtBuilder
	default:
		builder = srvr.jwtBuilder
		srvr.logger.Error().Str("token builder", srvr.config.TokenBuilder).Msgf("invalid token builder, using fallback builder")
	}

	payload, err := builder.VerifyToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.TokensVerifyResponse{
		Id:      payload.Id.String(),
		Iat:     timestamppb.New(payload.IssuedAt),
		Expires: timestamppb.New(payload.ExpiredAt),
		Payload: &rpc.TokensPayloadData{Data: payload.Payload.Data},
	}, nil
}
