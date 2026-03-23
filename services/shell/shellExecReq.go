package shell

import "github.com/ngyewch/fjage-go"

type ShellExecReq struct {
	fjage.Message

	Command    string   `json:"command,omitempty"`
	Script     string   `json:"script,omitempty"`
	ScriptArgs []string `json:"scriptArgs,omitempty"`
	Ans        bool     `json:"ans,omitempty"`
}

func (m ShellExecReq) Clazz() string {
	return "org.arl.fjage.shell.ShellExecReq"
}

func (m ShellExecReq) PropertiesMap() map[string]any {
	props := make(map[string]any)
	if m.Command != "" {
		props["command"] = m.Command
	}
	if m.Script != "" {
		props["script"] = m.Script
	}
	if m.ScriptArgs != nil {
		props["scriptArgs"] = m.ScriptArgs
	}
	props["ans"] = m.Ans
	return props
}
