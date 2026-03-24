package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doShellExec(ctx context.Context, cmd *cli.Command) error {
	command := cmd.StringArg(commandArg.Name)
	if command == "" {
		return fmt.Errorf("no command specified")
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellHelper, err := newShellHelper(ctx, cmd, gw)
			if err != nil {
				return err
			}

			err = shellHelper.ExecuteCommand(ctx, command)
			if err != nil {
				return err
			}

			return nil
		})
}
