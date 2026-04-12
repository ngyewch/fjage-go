package shell

import (
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/types"
)

type PutFileReq struct {
	*fjage.Message

	Filename string          `json:"filename"`
	Offset   int64           `json:"offset"`
	Contents types.ByteArray `json:"contents"`
}

func (m *PutFileReq) JavaClassName() string {
	return "org.arl.fjage.shell.PutFileReq"
}

func (m *PutFileReq) Header() *fjage.Message {
	return m.Message
}
