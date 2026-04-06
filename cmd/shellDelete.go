package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doShellDelete(ctx context.Context, cmd *cli.Command) error {
	remotePath := cmd.StringArg(remotePathArg.Name)

	if remotePath == "" {
		return fmt.Errorf("please specify %s", remotePathArg.Name)
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellHelper, err := newShellHelper(ctx, cmd, gw)
			if err != nil {
				return err
			}

			err = shellHelper.DeleteFile(ctx, remotePath)
			if err != nil {
				return err
			}

			return nil
		})
}
