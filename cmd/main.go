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

	serviceArg = &cli.StringArg{
		Name:      "service",
		UsageText: "(service)",
	}
	agentIdArg = &cli.StringArg{
		Name:      "agentID",
		UsageText: "(agentID)",
	}

	app = &cli.Command{
		Name:    "fjage-cli",
		Usage:   "fjage CLI",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "agents",
				Usage:  "agents",
				Action: doAgents,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
			},
			{
				Name:   "services",
				Usage:  "services",
				Action: doServices,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
			},
			{
				Name:   "agentForService",
				Usage:  "agent for service",
				Action: doAgentForService,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
				Arguments: []cli.Argument{
					serviceArg,
				},
			},
			{
				Name:   "agentsForService",
				Usage:  "agents for service",
				Action: doAgentsForService,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
				Arguments: []cli.Argument{
					serviceArg,
				},
			},
			{
				Name:   "containsAgent",
				Usage:  "contains agent",
				Action: doContainsAgent,
				Flags: []cli.Flag{
					gatewayUrlFlag,
				},
				Arguments: []cli.Argument{
					agentIdArg,
				},
			},
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
