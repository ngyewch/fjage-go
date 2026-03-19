package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	transport, err := newTransport(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(transport gateway.Transport) {
		_ = transport.Close()
	}(transport)

	gw, err := func(transport gateway.Transport) (gateway.Gateway, error) {
		return gateway.NewDefaultGateway(transport)
	}(transport)
	if err != nil {
		return err
	}
	defer func(gw gateway.Gateway) {
		_ = gw.Close()
	}(gw)

	{
		rsp, err := gw.Agents(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", rsp)
	}
	{
		rsp, err := gw.Services(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", rsp)
	}
	{
		rsp, err := gw.AgentForService(ctx, "org.arl.fjage.shell.Services.SHELL")
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", rsp)

		rsp1, err := gw.ContainsAgent(ctx, rsp.AgentID)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", rsp1)
	}
	{
		rsp, err := gw.AgentsForService(ctx, "org.arl.fjage.shell.Services.SHELL")
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", rsp)
	}

	return nil
}
