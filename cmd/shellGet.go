package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ngyewch/fjage-go/clients/shell"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doShellGet(ctx context.Context, cmd *cli.Command) error {
	remotePath := cmd.StringArg(remotePathArg.Name)
	localPath := cmd.StringArg(localPathArg.Name)

	if remotePath == "" {
		return fmt.Errorf("please specify %s", remotePathArg.Name)
	}
	if localPath == "" {
		return fmt.Errorf("please specify %s", localPathArg.Name)
	}

	_, err := os.Stat(filepath.Dir(localPath))
	if err != nil {
		return err
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellClient, err := shell.New(ctx, gw)
			if err != nil {
				return err
			}

			response, err := shellClient.GetFile(ctx, remotePath, 0, 0)
			if err != nil {
				return err
			}

			if response.Directory {
				return fmt.Errorf("%s is a directory", remotePath)
			}

			f, err := os.Create(localPath)
			if err != nil {
				return err
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			_, err = f.Write(response.Contents)
			if err != nil {
				return err
			}

			return nil
		})
}
