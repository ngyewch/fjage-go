package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doServices(ctx context.Context, cmd *cli.Command) error {
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			rsp, err := gw.Services(ctx)
			if err != nil {
				return err
			}
			for _, service := range rsp.Services {
				fmt.Println(service)
			}
			return nil
		})
}
