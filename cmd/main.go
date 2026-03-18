package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string

	gatewayUrlFlag = &cli.StringFlag{
		Name:    "gateway-url",
		Usage:   "gateway URL",
		Value:   "ws://localhost:8080/ws",
		Sources: cli.EnvVars("HTTP_PORT"),
	}

	app = &cli.Command{
		Name:    "fjage-cli",
		Usage:   "fjage CLI",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "test",
				Usage:  "test",
				Action: doTest,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
