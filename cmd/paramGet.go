package main

import (
	"context"

	"github.com/goforj/godump"
	"github.com/ngyewch/fjage-go/clients/param"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doParamGet(ctx context.Context, cmd *cli.Command) error {
	agentId := cmd.StringArg(agentIdArg.Name)
	timeout := cmd.Duration(timeoutFlag.Name)

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			paramHelper := param.NewParamHelper(gw)

			paramMap, err := paramHelper.GetParams(ctxWithTimeout, agentId)
			if err != nil {
				return err
			}

			godump.Dump(paramMap)

			return nil
		})
}
