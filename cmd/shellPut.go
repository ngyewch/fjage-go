package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

func doShellPut(ctx context.Context, cmd *cli.Command) error {
	copyBufferSize := cmd.Int(copyBufferSizeFlag.Name)
	progress := cmd.Bool(progressFlag.Name)
	localPath := cmd.StringArg(localPathArg.Name)
	remotePath := cmd.StringArg(remotePathArg.Name)

	if localPath == "" {
		return fmt.Errorf("please specify %s", localPathArg.Name)
	}
	if remotePath == "" {
		return fmt.Errorf("please specify %s", remotePathArg.Name)
	}

	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellClient, err := newShellClient(ctx, cmd, gw)
			if err != nil {
				return err
			}

			var pb *progressbar.ProgressBar
			if progress {
				pb = progressbar.DefaultBytes(fileInfo.Size())
				defer func(pb *progressbar.ProgressBar) {
					_ = pb.Finish()
					_ = pb.Close()
				}(pb)
			}

			if copyBufferSize <= 0 {
				contents, err := os.ReadFile(localPath)
				if err != nil {
					return err
				}

				err = shellClient.PutFile(ctx, remotePath, 0, contents)
				if err != nil {
					return err
				}

				if pb != nil {
					err = pb.Add(len(contents))
					if err != nil {
						return err
					}
				}
			} else {
				f, err := os.Open(localPath)
				if err != nil {
					return err
				}
				defer func(f *os.File) {
					_ = f.Close()
				}(f)

				var offset int64 = 0
				buffer := make([]byte, copyBufferSize)
				for {
					readLen, err := f.Read(buffer)
					if err != nil {
						if err == io.EOF {
							break
						} else {
							return err
						}
					}

					err = shellClient.PutFile(ctx, remotePath, offset, buffer[0:readLen])
					if err != nil {
						return err
					}

					offset += int64(readLen)
					if pb != nil {
						err = pb.Add(readLen)
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		})
}
