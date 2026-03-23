package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

func doShellGet(ctx context.Context, cmd *cli.Command) error {
	progress := cmd.Bool(progressFlag.Name)
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

	var progressHandler func(current int64, total int64)
	var pb *progressbar.ProgressBar
	if progress {
		progressHandler = func(current int64, total int64) {
			if pb == nil {
				pb = progressbar.DefaultBytes(total)
			}
			_ = pb.Set64(current)
		}
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellHelper, err := newShellHelper(ctx, cmd, gw)
			if err != nil {
				return err
			}

			err = shellHelper.GetFile(ctx, remotePath, localPath, progressHandler)
			if err != nil {
				return err
			}

			return nil
		})
}
