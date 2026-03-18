package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func newGateway(ctx context.Context, cmd *cli.Command) (gateway.Gateway, error) {
	gatewayUrl := cmd.String(gatewayUrlFlag.Name)

	u, err := url.Parse(gatewayUrl)
	if err != nil {
		return nil, err
	}

	if (u.Scheme == "ws") || (u.Scheme == "wss") {
		return gateway.NewWebSocketGateway(ctx, gatewayUrl)
	}

	return nil, fmt.Errorf("unsupported gateway url scheme: %s", u.Scheme)
}
