package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jwalton/go-supportscolor"
	"github.com/phsym/console-slog"
	"github.com/urfave/cli/v3"
)

var (
	version string

	gatewayUrlFlag = &cli.StringFlag{
		Name:    "gateway-url",
		Usage:   "gateway URL",
		Value:   "ws://localhost:8080/ws",
		Sources: cli.EnvVars("GATEWAY_URL"),
	}
	copyBufferSizeFlag = &cli.IntFlag{
		Name:    "copy-buffer-size",
		Usage:   "copy buffer size (bytes)",
		Value:   16 * 1024,
		Sources: cli.EnvVars("COPY_BUFFER_SIZE"),
	}
	progressFlag = &cli.BoolFlag{
		Name:    "progress",
		Usage:   "show progress",
		Value:   true,
		Sources: cli.EnvVars("PROGRESS"),
	}

	serviceArg = &cli.StringArg{
		Name:      "service",
		UsageText: "(service)",
	}
	agentIdArg = &cli.StringArg{
		Name:      "agentID",
		UsageText: "(agentID)",
	}
	pathArg = &cli.StringArg{
		Name:      "path",
		UsageText: "(path)",
	}
	localPathArg = &cli.StringArg{
		Name:      "local-path",
		UsageText: "(local path)",
	}
	remotePathArg = &cli.StringArg{
		Name:      "remote-path",
		UsageText: "(remote path)",
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
				Name:  "shell",
				Usage: "shell",
				Commands: []*cli.Command{
					{
						Name:   "ls",
						Usage:  "ls",
						Action: doShellLs,
						Flags: []cli.Flag{
							gatewayUrlFlag,
						},
						Arguments: []cli.Argument{
							pathArg,
						},
					},
					{
						Name:   "get",
						Usage:  "get",
						Action: doShellGet,
						Flags: []cli.Flag{
							gatewayUrlFlag,
							copyBufferSizeFlag,
							progressFlag,
						},
						Arguments: []cli.Argument{
							remotePathArg,
							localPathArg,
						},
					},
					{
						Name:   "put",
						Usage:  "put",
						Action: doShellPut,
						Flags: []cli.Flag{
							gatewayUrlFlag,
							copyBufferSizeFlag,
							progressFlag,
						},
						Arguments: []cli.Argument{
							localPathArg,
							remotePathArg,
						},
					},
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

func init() {
	logLevel := slog.LevelInfo
	_ = logLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	addSource := os.Getenv("LOG_ADD_SOURCE") == "true"
	noColor := os.Getenv("NO_COLOR") == "true"

	var logger *slog.Logger
	if supportscolor.Stderr().SupportsColor {
		logger = slog.New(
			console.NewHandler(os.Stderr, &console.HandlerOptions{
				Level:     logLevel,
				AddSource: addSource,
				NoColor:   noColor,
			}),
		)
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: addSource,
		}))
	}
	slog.SetDefault(logger)
}
