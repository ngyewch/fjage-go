package main

import (
	"context"
	"time"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	gw, err := newGateway(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(gw gateway.Gateway) {
		_ = gw.Close()
	}(gw)

	err = gw.Agents(ctx)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	/*
		_ = gw.Close()

		time.Sleep(2 * time.Second)

	*/

	return nil
}
