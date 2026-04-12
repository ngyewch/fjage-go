package shell

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/ngyewch/fjage-go/services/shell"
	"github.com/ngyewch/fjage-go/types"
)

type Client struct {
	gw           gateway.Gateway
	shellAgentID string
}

func New(ctx context.Context, gw gateway.Gateway) (*Client, error) {
	getAgentForServiceResponse, err := gw.AgentForService(ctx, shell.ServiceName)
	if err != nil {
		return nil, err
	}
	shellAgentID := getAgentForServiceResponse.AgentID
	if shellAgentID == "" {
		return nil, fmt.Errorf("could not find agent for %s", shell.ServiceName)
	}
	return &Client{
		gw:           gw,
		shellAgentID: shellAgentID,
	}, nil
}

func (client *Client) Gateway() gateway.Gateway {
	return client.gw
}

func (client *Client) GetFile(ctx context.Context, filename string, offset int64, length int64) (*shell.GetFileRsp, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	request := &shell.GetFileReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    client.shellAgentID,
			Sender:       client.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Filename: filename,
		Offset:   offset,
		Length:   length,
	}
	sendResponse, err := client.gw.Send(ctx, request)
	if err != nil {
		return nil, err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return nil, fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	switch m := sendResponse.Message.(type) {
	case *shell.GetFileRsp:
		return m, nil
	default:
		return nil, nil
	}
}

func (client *Client) PutFile(ctx context.Context, filename string, offset int64, contents []byte) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	request := &shell.PutFileReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    client.shellAgentID,
			Sender:       client.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Filename: filename,
		Offset:   offset,
		Contents: types.ByteArray(contents),
	}
	sendResponse, err := client.gw.Send(ctx, request)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeAgree {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	return nil
}

func (client *Client) ExecuteCommand(ctx context.Context, command string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	request := &shell.ShellExecReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    client.shellAgentID,
			Sender:       client.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Command: command,
		Ans:     true,
	}
	sendResponse, err := client.gw.Send(ctx, request)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeAgree {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	return nil
}

func (client *Client) ExecuteScript(ctx context.Context, scriptFile string, scriptArgs []string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	request := &shell.ShellExecReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    client.shellAgentID,
			Sender:       client.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Script:     scriptFile,
		ScriptArgs: scriptArgs,
		Ans:        true,
	}
	sendResponse, err := client.gw.Send(ctx, request)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeAgree {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	return nil
}
