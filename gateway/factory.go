package gateway

import (
	"encoding/json"

	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/param"
	"github.com/ngyewch/fjage-go/services/shell"
)

func instantiateMessage(clazz string) fjage.IMessage {
	switch clazz {
	case "org.arl.fjage.param.ParameterReq":
		return new(param.ParameterReq)
	case "org.arl.fjage.param.ParameterRsp":
		return new(param.ParameterRsp)
	case "org.arl.fjage.shell.GetFileReq":
		return new(shell.GetFileReq)
	case "org.arl.fjage.shell.GetFileRsp":
		return new(shell.GetFileRsp)
	case "org.arl.fjage.shell.PutFileReq":
		return new(shell.PutFileReq)
	case "org.arl.fjage.shell.ShellExecReq":
		return new(shell.ShellExecReq)
	case "org.arl.fjage.Message":
		return new(fjage.Message)
	default:
		return nil
	}
}

func unmarshalMessage(messageEnvelope *MessageEnvelope) (fjage.IMessage, error) {
	if messageEnvelope == nil {
		return nil, nil
	}
	jsonBytes, err := json.Marshal(messageEnvelope.Data)
	if err != nil {
		return nil, err
	}
	m := instantiateMessage(messageEnvelope.Clazz)
	if m == nil {
		m = new(fjage.Message)
	}
	if m == nil {
		return nil, nil
	}
	err = json.Unmarshal(jsonBytes, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
