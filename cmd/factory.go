package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ngyewch/fjage-go/clients/shell"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func newTransport(ctx context.Context, cmd *cli.Command) (gateway.Transport, error) {
	gatewayUrl := cmd.String(gatewayUrlFlag.Name)

	u, err := url.Parse(gatewayUrl)
	if err != nil {
		return nil, err
	}

	if (u.Scheme == "ws") || (u.Scheme == "wss") {
		return gateway.NewWebSocketTransport(ctx, gatewayUrl)
	}

	return nil, fmt.Errorf("unsupported gateway url scheme: %s", u.Scheme)
}

func withGateway(ctx context.Context, cmd *cli.Command, handler func(gw gateway.Gateway) error) error {
	transport, err := newTransport(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(transport gateway.Transport) {
		_ = transport.Close()
	}(transport)

	gw, err := func(transport gateway.Transport) (gateway.Gateway, error) {
		return gateway.NewDefaultGateway(ctx, transport)
	}(transport)
	if err != nil {
		return err
	}
	defer func(gw gateway.Gateway) {
		_ = gw.Close()
	}(gw)

	return handler(gw)
}

func newShellClient(ctx context.Context, cmd *cli.Command, gw gateway.Gateway) (*shell.Client, error) {
	return shell.New(ctx, gw)
}

func newShellHelper(ctx context.Context, cmd *cli.Command, gw gateway.Gateway) (*shell.Helper, error) {
	copyBufferSize := cmd.Int(copyBufferSizeFlag.Name)

	client, err := newShellClient(ctx, cmd, gw)
	if err != nil {
		return nil, err
	}

	helperOptions := shell.DefaultHelperOptions
	helperOptions.CopyBufferSize = int64(copyBufferSize)
	return shell.NewHelper(client, &helperOptions)
}
