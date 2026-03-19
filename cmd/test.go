package main

import (
	"context"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			return nil
		})
}
