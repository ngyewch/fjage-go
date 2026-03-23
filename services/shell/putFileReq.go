package shell

import (
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/types"
)

type PutFileReq struct {
	fjage.Message

	Filename string          `json:"filename"`
	Offset   int64           `json:"offset"`
	Contents types.ByteArray `json:"contents"`
}

func (m PutFileReq) Clazz() string {
	return "org.arl.fjage.shell.PutFileReq"
}

func (m PutFileReq) PropertiesMap() map[string]any {
	props := make(map[string]any)
	if m.Filename != "" {
		props["filename"] = m.Filename
	}
	props["offset"] = m.Offset
	if m.Contents != nil {
		props["contents"] = m.Contents
	}
	return props
}
