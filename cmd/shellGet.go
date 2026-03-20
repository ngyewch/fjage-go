package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/ngyewch/fjage-go/services/shell"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

func doShellGet(ctx context.Context, cmd *cli.Command) error {
	copyBufferSize := cmd.Int(copyBufferSizeFlag.Name)
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

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellClient, err := newShellClient(ctx, cmd, gw)
			if err != nil {
				return err
			}

			remoteDir := path.Dir(remotePath)
			if remoteDir == "" {
				remoteDir = "."
			}
			var dirEntry *shell.DirEntry
			{
				getFileRsp, err := shellClient.GetFile(ctx, remoteDir, 0, 0)
				if err != nil {
					return err
				}
				dirEntries, err := getFileRsp.DirEntries()
				if err != nil {
					return err
				}
				remoteBase := path.Base(remotePath)
				for _, entry := range dirEntries {
					if entry.Name == remoteBase {
						dirEntry = &entry
						break
					}
				}
			}
			if dirEntry == nil {
				return fmt.Errorf("remote file not found")
			}

			f, err := os.Create(localPath)
			if err != nil {
				return err
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			var pb *progressbar.ProgressBar
			if progress {
				pb = progressbar.DefaultBytes(dirEntry.Size)
				defer func(pb *progressbar.ProgressBar) {
					_ = pb.Finish()
					_ = pb.Close()
				}(pb)
			}

			if copyBufferSize <= 0 {
				response, err := shellClient.GetFile(ctx, remotePath, 0, 0)
				if err != nil {
					return err
				}

				if response.Directory {
					return fmt.Errorf("%s is a directory", remotePath)
				}

				_, err = f.Write(response.Contents)
				if err != nil {
					return err
				}

				if pb != nil {
					err = pb.Add(len(response.Contents))
					if err != nil {
						return err
					}
				}
			} else {
				var offset int64 = 0

				for offset < dirEntry.Size {
					length := int64(copyBufferSize)
					if offset+length > dirEntry.Size {
						length = dirEntry.Size - offset
					}
					response, err := shellClient.GetFile(ctx, remotePath, offset, length)
					if err != nil {
						return err
					}

					if response.Directory {
						return fmt.Errorf("%s is a directory", remotePath)
					}

					_, err = f.Write(response.Contents)
					if err != nil {
						return err
					}

					if pb != nil {
						err = pb.Add(len(response.Contents))
						if err != nil {
							return err
						}
					}

					offset += length
				}
			}

			return nil
		})
}
