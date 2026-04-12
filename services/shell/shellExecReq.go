package shell

import "github.com/ngyewch/fjage-go"

type ShellExecReq struct {
	*fjage.Message

	Command    string   `json:"command,omitempty"`
	Script     string   `json:"script,omitempty"`
	ScriptArgs []string `json:"scriptArgs,omitempty"`
	Ans        bool     `json:"ans,omitempty"`
}

func (m *ShellExecReq) JavaClassName() string {
	return "org.arl.fjage.shell.ShellExecReq"
}

func (m *ShellExecReq) Header() *fjage.Message {
	return m.Message
}
