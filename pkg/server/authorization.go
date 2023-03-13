package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

func (srvr *RPCTokensServer) authorize(ctx context.Context) error {
	// If there is no service secret then there is no authorize
	// Error messsages are already logged in main.go
	if srvr.config.ApiSecret == "" {
		return nil
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing metadata")
	}
	values := meta.Get(srvr.config.ApiHeader)
	if len(values) == 0 {
		return fmt.Errorf("missing %s ", srvr.config.ApiHeader)
	}

	secret := values[0]
	if secret == srvr.config.ApiSecret {
		return nil
	}

	return fmt.Errorf("invalid %s", srvr.config.ApiHeader)
}
