package shell

import "github.com/ngyewch/fjage-go"

type GetFileReq struct {
	fjage.Message

	Filename string `json:"filename"`
	Offset   int64  `json:"offset"`
	Length   int64  `json:"length"`
}

func (m GetFileReq) Clazz() string {
	return "org.arl.fjage.shell.GetFileReq"
}

func (m GetFileReq) PropertiesMap() map[string]any {
	return map[string]any{
		"filename": m.Filename,
		"offset":   m.Offset,
		"length":   m.Length,
	}
}
