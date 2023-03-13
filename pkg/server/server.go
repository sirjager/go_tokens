package server

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sirjager/go_tokens/cfg"
	"github.com/sirjager/go_tokens/pkg/tokens"

	rpc "github.com/sirjager/rpcs/tokens/go"
)

type RPCTokensServer struct {
	rpc.UnimplementedTokensServer
	startTime     time.Time
	config        cfg.Config
	logger        zerolog.Logger
	pasetoBuilder tokens.TokenBuilder
	jwtBuilder    tokens.TokenBuilder
}

func NewTokensServer(startTime time.Time, logger zerolog.Logger, config cfg.Config) (*RPCTokensServer, error) {
	pasetoBuilder, err := tokens.NewPasetoBuilder(config.TokenSecret)
	if err != nil {
		logger.Error().Err(err).Msg("unable to create token builder")
		return nil, err
	}
	jwtBuilder, err := tokens.NewJWTBuilder(config.TokenSecret)
	if err != nil {
		logger.Error().Err(err).Msg("unable to create token builder")
		return nil, err
	}

	srvic := &RPCTokensServer{
		startTime:     startTime,
		config:        config,
		logger:        logger,
		jwtBuilder:    jwtBuilder,
		pasetoBuilder: pasetoBuilder,
	}
	return srvic, nil
}
