package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doContainsAgent(ctx context.Context, cmd *cli.Command) error {
	agentID := cmd.StringArg(agentIdArg.Name)
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			rsp, err := gw.ContainsAgent(ctx, agentID)
			if err != nil {
				return err
			}
			fmt.Println(rsp.Answer)
			return nil
		})
}
