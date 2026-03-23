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
	props := make(map[string]any)
	if m.Filename != "" {
		props["filename"] = m.Filename
	}
	props["offset"] = m.Offset
	props["length"] = m.Length
	return props
}
