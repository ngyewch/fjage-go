package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/gateway"
)

const (
	serviceName = "org.arl.fjage.shell.Services.SHELL"
)

type Shell struct {
	gw           gateway.Gateway
	shellAgentID string
}

func NewShell(ctx context.Context, gw gateway.Gateway) (*Shell, error) {
	getAgentForServiceResponse, err := gw.AgentForService(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	shellAgentID := getAgentForServiceResponse.AgentID
	if shellAgentID == "" {
		return nil, fmt.Errorf("could not find agent for %s", serviceName)
	}
	return &Shell{
		gw:           gw,
		shellAgentID: shellAgentID,
	}, nil
}

func (shell *Shell) GetFile(ctx context.Context, filename string, offset int64, length int64) (*GetFileRsp, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	request := &GetFileReq{
		Message: fjage.Message{
			MsgID:     id.String(),
			Perf:      "REQUEST",
			Recipient: shell.shellAgentID,
			Sender:    shell.gw.AgentID(),
			SentAt:    time.Now().UnixMilli(),
		},
		Filename: filename,
		Offset:   offset,
		Length:   length,
	}
	sendResponse, err := shell.gw.Send(ctx, request.Clazz(), &request.Message, request.PropertiesMap())
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(sendResponse.Message)
	if err != nil {
		return nil, err
	}

	var messageWrapper gateway.MessageWrapper[*GetFileRsp]
	err = json.Unmarshal(jsonBytes, &messageWrapper)
	if err != nil {
		return nil, err
	}
	if messageWrapper.Data.Perf != "INFORM" {
		return nil, fjage.NewPerformativeError(messageWrapper.Data.Perf)
	}

	return messageWrapper.Data, nil
}

func (shell *Shell) PutFile(ctx context.Context, filename string, offset int64, contents []byte) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	request := &PutFileReq{
		Message: fjage.Message{
			MsgID:     id.String(),
			Perf:      "REQUEST",
			Recipient: shell.shellAgentID,
			Sender:    shell.gw.AgentID(),
			SentAt:    time.Now().UnixMilli(),
		},
		Filename: filename,
		Offset:   offset,
		Contents: &gateway.Array{
			Clazz: "[B",
			Data:  contents,
		},
	}
	sendResponse, err := shell.gw.Send(ctx, request.Clazz(), &request.Message, request.PropertiesMap())
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(sendResponse.Message)
	if err != nil {
		return err
	}

	var messageWrapper gateway.MessageWrapper[*fjage.Message]
	err = json.Unmarshal(jsonBytes, &messageWrapper)
	if err != nil {
		return err
	}
	if messageWrapper.Data.Perf != "AGREE" {
		return fjage.NewPerformativeError(messageWrapper.Data.Perf)
	}

	return nil
}
