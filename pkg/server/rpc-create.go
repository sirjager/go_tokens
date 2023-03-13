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

func (srvr *RPCTokensServer) TokensCreate(ctx context.Context, req *rpc.TokensCreateRequest) (*rpc.TokensCreateResponse, error) {
	// This only validates if the request is authorized or not
	if err := srvr.authorize(ctx); err != nil {
		return nil, unAuthenticatedError(err)
	}

	// Since duration is optional, so by default we will set alive duration to 10 mins
	duration := srvr.config.DefaultTokenDuration // time.Second *  seconds * mins * hrs * days

	if duration.Seconds() > srvr.config.MaximumTokenDuration.Seconds() {
		return nil, status.Errorf(codes.Internal, "tokens duration exceeds maximum limit set be server")
	}

	var builder tokens.TokenBuilder
	switch strings.ToLower(srvr.config.TokenBuilder) {
	case tokens.Paseto:
		builder = srvr.pasetoBuilder
	case tokens.Jwt:
		builder = srvr.jwtBuilder
	default:
		srvr.logger.Error().Str("token builder", srvr.config.TokenBuilder).Msgf("invalid token builder, using fallback builder")
		builder = srvr.jwtBuilder
	}

	access_token, payload, err := builder.CreateToken(tokens.PayloadData{Data: req.Payload}, duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.TokensCreateResponse{
		Token: access_token,
		Payload: &rpc.TokensPayload{
			Id:      payload.Id.String(),
			Iat:     timestamppb.New(payload.IssuedAt),
			Expires: timestamppb.New(payload.ExpiredAt),
			Payload: &rpc.TokensPayloadData{Data: payload.Payload.Data},
		},
	}, nil
}
