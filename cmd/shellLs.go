package main

import (
	"context"
	"os"
	"time"

	"github.com/ngyewch/fjage-go/clients/shell"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/urfave/cli/v3"
)

func doShellLs(ctx context.Context, cmd *cli.Command) error {
	path := cmd.StringArg(pathArg.Name)
	if path == "" {
		path = "."
	}

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			shellClient, err := shell.New(ctx, gw)
			if err != nil {
				return err
			}

			response, err := shellClient.GetFile(ctx, path, 0, 0)
			if err != nil {
				return err
			}

			dirEntries, err := response.DirEntries()
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.Header("Name", "Size", "Last Modified")
			table.Configure(func(config *tablewriter.Config) {
				config.Row.Alignment.PerColumn = []tw.Align{tw.AlignLeft, tw.AlignRight, tw.AlignLeft}
			})
			for _, entry := range dirEntries {
				err = table.Append([]any{
					entry.Name,
					entry.Size,
					entry.LastModified.Format(time.RFC3339),
				})
				if err != nil {
					return err
				}
			}
			err = table.Render()
			if err != nil {
				return err
			}

			return nil
		})
}
