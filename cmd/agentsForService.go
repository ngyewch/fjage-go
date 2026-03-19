package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doAgentsForService(ctx context.Context, cmd *cli.Command) error {
	service := cmd.StringArg(serviceArg.Name)
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			rsp, err := gw.AgentsForService(ctx, service)
			if err != nil {
				return err
			}
			for _, agentID := range rsp.AgentIDs {
				fmt.Println(agentID)
			}
			return nil
		})
}
