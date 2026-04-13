package param

import "github.com/ngyewch/fjage-go"

type ParameterRsp struct {
	*fjage.Message

	Index    int                      `json:"index"`
	Param    string                   `json:"param,omitempty"`
	Readonly []string                 `json:"readonly,omitempty"`
	Value    *GenericValue            `json:"value,omitempty"`
	Values   map[string]*GenericValue `json:"values,omitempty"`
}

func (m *ParameterRsp) JavaClassName() string {
	return "org.arl.fjage.param.ParameterRsp"
}

func (m *ParameterRsp) Header() *fjage.Message {
	return m.Message
}
