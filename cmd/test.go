package main

import (
	"context"
	"fmt"
	"time"

	"github.com/goforj/godump"
	paramClient "github.com/ngyewch/fjage-go/clients/param"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/urfave/cli/v3"
	"go.octolab.org/pointer"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	return withGateway(ctx, cmd,
		func(gw gateway.Gateway) error {
			paramHelper := paramClient.NewParamHelper(gw)
			{
				agentForServiceResponse, err := gw.AgentForService(ctx, "org.arl.fjage.shell.Services.SHELL")
				if err != nil {
					return err
				}
				shellAgentID := agentForServiceResponse.AgentID

				var language string
				err = paramHelper.GetParam(ctx, shellAgentID, "org.arl.fjage.shell.ShellParam.language", &language)
				if err != nil {
					return err
				}
				fmt.Printf("language: %s\n", language)
			}
			{
				var strings []string
				err := paramHelper.GetParam(ctx, "test", "TestParam.strings", &strings)
				if err != nil {
					return err
				}
				fmt.Printf("strings: %v\n", strings)
			}
			{
				t := time.Now().UnixMilli()
				err := paramHelper.SetParam(ctx, "test", "TestParam.strings", []string{
					fmt.Sprintf("boo-%d", t),
					fmt.Sprintf("hoo-%d", t),
				})
				if err != nil {
					return err
				}
			}
			{
				var strings []string
				err := paramHelper.GetParam(ctx, "test", "TestParam.strings", &strings)
				if err != nil {
					return err
				}
				fmt.Printf("strings: %v\n", strings)
			}
			{
				data, err := paramHelper.GetParams(ctx, "test")
				if err != nil {
					return err
				}
				godump.Dump(data)
			}
			{
				type MyData struct {
					Strings []string `fjage.parameter:"TestParam.strings"`
					Ints    []int32  `fjage.parameter:"TestParam.ints"`
					String1 string   `fjage.parameter:"TestParam.string1"`
				}
				var data MyData
				err := paramHelper.GetParamsAndPopulate(ctx, "test", &data)
				if err != nil {
					return err
				}
				godump.Dump(data)
			}
			{
				type MyData struct {
					Strings []string `fjage.parameter:"TestParam.strings"`
					Ints    []int32  `fjage.parameter:"TestParam.ints"`
					String1 *string  `fjage.parameter:"TestParam.string1"`
				}
				var data MyData
				err := paramHelper.GetParamsAndPopulate(ctx, "test", &data)
				if err != nil {
					return err
				}
				godump.Dump(data)
			}
			{
				type MyData struct {
					String1 *string  `fjage.parameter:"TestParam.string1"`
					Int1    *int32   `fjage.parameter:"TestParam.int1"`
					Float1  *float32 `fjage.parameter:"TestParam.float1"`
					Double1 *float64 `fjage.parameter:"TestParam.double1"`
				}
				{
					var data MyData
					err := paramHelper.GetParamsAndPopulate(ctx, "test", &data)
					if err != nil {
						return err
					}
					godump.Dump(data)
				}
				{
					err := paramHelper.SetParamsFromStruct(ctx, "test", MyData{
						String1: pointer.ToString(time.Now().Format(time.RFC3339)),
					})
					if err != nil {
						return err
					}
					var data MyData
					err = paramHelper.GetParamsAndPopulate(ctx, "test", &data)
					if err != nil {
						return err
					}
					godump.Dump(data)
				}
			}
			return nil
		})
}
