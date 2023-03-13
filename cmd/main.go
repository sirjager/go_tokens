package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"

	"github.com/sirjager/go_tokens/cfg"

	"github.com/sirjager/go_tokens/cmd/server"
	tokensServer "github.com/sirjager/go_tokens/pkg/server"
)

var logger zerolog.Logger
var startTime time.Time

const serviceName = "tokens"

func init() {
	startTime = time.Now()
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = logger.With().Timestamp().Logger()
	logger = logger.With().Str("service", serviceName).Logger()
}

func main() {
	config, err := cfg.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to load configs")
	}

	srvr, err := tokensServer.NewTokensServer(startTime, logger, config)
	if err != nil {
		logger.Fatal().Err(err).Msgf("unable to create %s service", serviceName)
	}

	errs := make(chan error)

	if config.ApiSecret == "" {
		logger.Warn().Str(config.ApiHeader, "missing").Msg(color.RedString("service is unprotected"))
	}

	go handleSignals(errs)
	go server.RunGRPCServer(srvr, logger, config, errs)

	if config.Http != "" {
		// Made this optional since this service will only be used by internal services
		go server.RunGatewayServer(srvr, logger, config, errs)
	}

	logger.Error().Err(<-errs).Msg("exit")
}

func handleSignals(errs chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
}
