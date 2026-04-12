package main

import (
	"context"

	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			/*
				{
					agentForServiceResponse, err := gw.AgentForService(ctx, "org.arl.unet.Services.DEVICE_INFO")
					if err != nil {
						return err
					}
					fmt.Printf("agentID = %s\n", agentForServiceResponse.AgentID)
					deviceInfoAgentID := agentForServiceResponse.AgentID

					id, err := uuid.NewRandom()
					if err != nil {
						return err
					}
					sendResponse, err := gw.Send(ctx, "org.arl.fjage.param.ParameterReq", &fjage.Message{
						MsgID:     id.String(),
						Perf:      "REQUEST",
						Recipient: deviceInfoAgentID,
						Sender:    gw.AgentID(),
						SentAt:    time.Now().UnixMilli(),
					}, map[string]any{
						"param": "model",
					})
					if err != nil {
						return err
					}

					fmt.Printf("sendResponse = %+v\n", sendResponse)
					godump.Dump(sendResponse)
				}
				{
					agentForServiceResponse, err := gw.AgentForService(ctx, "org.arl.unet.Services.BASEBAND")
					if err != nil {
						return err
					}
					fmt.Printf("agentID = %s\n", agentForServiceResponse.AgentID)
					deviceInfoAgentID := agentForServiceResponse.AgentID

					id, err := uuid.NewRandom()
					if err != nil {
						return err
					}
					sendResponse, err := gw.Send(ctx, "org.arl.fjage.param.ParameterReq", &fjage.Message{
						MsgID:     id.String(),
						Perf:      "REQUEST",
						Recipient: deviceInfoAgentID,
						Sender:    gw.AgentID(),
						SentAt:    time.Now().UnixMilli(),
					}, map[string]any{
						"param": "record",
					})
					if err != nil {
						return err
					}

					fmt.Printf("sendResponse = %+v\n", sendResponse)
					godump.Dump(sendResponse)
				}
			*/

			return nil
		})
}
