package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/sirjager/go_tokens/cfg"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/sirjager/go_tokens/docs/statik"
	"github.com/sirjager/go_tokens/pkg/server"

	rpc "github.com/sirjager/rpcs/tokens/go"
)

func RunGatewayServer(srvr *server.RPCTokensServer, logger zerolog.Logger, config cfg.Config, errs chan error) {
	opts := []runtime.ServeMuxOption{}

	opts = append(opts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	}))

	opts = append(opts, runtime.WithIncomingHeaderMatcher(AllowedHeaders([]string{config.ApiHeader})))

	grpcMux := runtime.NewServeMux(opts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rpc.RegisterTokensHandlerServer(ctx, grpcMux, srvr)
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// File server for swagger documentations
	statikFS, err := fs.New()
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHander := http.StripPrefix("/api/swagger/", http.FileServer(statikFS))
	mux.Handle("/api/swagger/", swaggerHander)

	mux.Handle("/metrics", promhttp.Handler())

	listener, err := net.Listen("tcp", ":"+config.Http)
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("unable to listen grpc tcp server")
	}

	logger.Info().Msgf("started HTTP server at %s", listener.Addr().String())

	handler := HTTPLogger(logger, mux)

	errs <- http.Serve(listener, handler)
}
