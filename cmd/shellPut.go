package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/ngyewch/fjage-go/services/shell"
	"github.com/urfave/cli/v3"
)

func doShellPut(ctx context.Context, cmd *cli.Command) error {
	localPath := cmd.StringArg(localPathArg.Name)
	remotePath := cmd.StringArg(remotePathArg.Name)

	if localPath == "" {
		return fmt.Errorf("please specify %s", localPathArg.Name)
	}
	if remotePath == "" {
		return fmt.Errorf("please specify %s", remotePathArg.Name)
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellClient, err := shell.NewShell(ctx, gw)
			if err != nil {
				return err
			}

			contents, err := os.ReadFile(localPath)
			if err != nil {
				return err
			}

			err = shellClient.PutFile(ctx, remotePath, 0, contents)
			if err != nil {
				return err
			}

			return nil
		})
}
