package shell

import "github.com/ngyewch/fjage-go"

type GetFileReq struct {
	*fjage.Message

	Filename string `json:"filename"`
	Offset   int64  `json:"offset"`
	Length   int64  `json:"length"`
}

func (m *GetFileReq) JavaClassName() string {
	return "org.arl.fjage.shell.GetFileReq"
}

func (m *GetFileReq) Header() *fjage.Message {
	return m.Message
}
