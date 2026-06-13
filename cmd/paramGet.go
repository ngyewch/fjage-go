package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ngyewch/fjage-go/clients/param"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doParamGet(ctx context.Context, cmd *cli.Command) error {
	agentId := cmd.StringArg(agentIdArg.Name)
	timeout := cmd.Duration(timeoutFlag.Name)

	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			paramHelper := param.NewParamHelper(gw)

			paramMap := make(map[string]any)
			err := paramHelper.GetParams(ctxWithTimeout, agentId, paramMap)
			if err != nil {
				return err
			}

			for key, value := range paramMap {
				jsonBytes, err := json.Marshal(value)
				if err != nil {
					return err
				}
				valueAsJson := string(jsonBytes)
				fmt.Printf("%s = %s (%T)\n", key, valueAsJson, value)
			}

			return nil
		})
}
