package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doAgents(ctx context.Context, cmd *cli.Command) error {
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			rsp, err := gw.Agents(ctx)
			if err != nil {
				return err
			}
			for _, agentID := range rsp.AgentIDs {
				fmt.Println(agentID)
			}
			return nil
		})
}
