package gateway

import (
	"encoding/json"

	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/param"
	"github.com/ngyewch/fjage-go/services/shell"
)

type MessageFactory struct {
	instantiatorMap map[string]func() fjage.IMessage
}

func NewMessageFactory() *MessageFactory {
	instantiatorMap := map[string]func() fjage.IMessage{
		"org.arl.fjage.param.ParameterReq": func() fjage.IMessage {
			return new(param.ParameterReq)
		},
		"org.arl.fjage.param.ParameterRsp": func() fjage.IMessage {
			return new(param.ParameterRsp)
		},
		"org.arl.fjage.shell.GetFileReq": func() fjage.IMessage {
			return new(shell.GetFileReq)
		},
		"org.arl.fjage.shell.GetFileRsp": func() fjage.IMessage {
			return new(shell.GetFileRsp)
		},
		"org.arl.fjage.shell.PutFileReq": func() fjage.IMessage {
			return new(shell.PutFileReq)
		},
		"org.arl.fjage.shell.ShellExecReq": func() fjage.IMessage {
			return new(shell.ShellExecReq)
		},
		"org.arl.fjage.Message": func() fjage.IMessage {
			return new(fjage.Message)
		},
	}
	return &MessageFactory{
		instantiatorMap: instantiatorMap,
	}
}

func (messageFactory *MessageFactory) Register(clazz string, instantiator func() fjage.IMessage) {
	messageFactory.instantiatorMap[clazz] = instantiator
}

func (messageFactory *MessageFactory) InstantiateMessage(clazz string) fjage.IMessage {
	instantiator, ok := messageFactory.instantiatorMap[clazz]
	if ok {
		return instantiator()
	}
	return nil
}

func (messageFactory *MessageFactory) UnmarshalMessage(messageEnvelope *MessageEnvelope) (fjage.IMessage, error) {
	if messageEnvelope == nil {
		return nil, nil
	}
	jsonBytes, err := json.Marshal(messageEnvelope.Data)
	if err != nil {
		return nil, err
	}
	m := messageFactory.InstantiateMessage(messageEnvelope.Clazz)
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
