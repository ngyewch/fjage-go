package param

import "github.com/ngyewch/fjage-go"

type ParameterReq struct {
	*fjage.Message

	Param    string              `json:"param,omitempty"`
	Requests []ParameterReqEntry `json:"requests,omitempty"`
	Value    *GenericValue       `json:"value,omitempty"`
}

func (m *ParameterReq) JavaClassName() string {
	return "org.arl.fjage.param.ParameterReq"
}

func (m *ParameterReq) Header() *fjage.Message {
	return m.Message
}

type ParameterReqEntry struct {
	Param string        `json:"param,omitempty"`
	Value *GenericValue `json:"value,omitempty"`
}
