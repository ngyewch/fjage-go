package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doAgentForService(ctx context.Context, cmd *cli.Command) error {
	service := cmd.StringArg(serviceArg.Name)
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			rsp, err := gw.AgentForService(ctx, service)
			if err != nil {
				return err
			}
			fmt.Println(rsp.AgentID)
			return nil
		})
}
