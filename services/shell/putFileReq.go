package shell

import (
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/types"
)

type PutFileReq struct {
	fjage.Message

	Filename string          `json:"filename"`
	Offset   int64           `json:"ofs"` // TODO offset
	Contents types.ByteArray `json:"contents"`
}

func (m PutFileReq) Clazz() string {
	return "org.arl.fjage.shell.PutFileReq"
}

func (m PutFileReq) PropertiesMap() map[string]any {
	return map[string]any{
		"filename": m.Filename,
		"ofs":      m.Offset, // TODO offset
		"contents": m.Contents,
	}
}
