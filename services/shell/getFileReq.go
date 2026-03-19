package shell

import "github.com/ngyewch/fjage-go"

type GetFileReq struct {
	fjage.Message

	Filename string `json:"filename"`
	Offset   int64  `json:"ofs"` // TODO offset
	Length   int64  `json:"len"` // TODO length
}

func (m GetFileReq) Clazz() string {
	return "org.arl.fjage.shell.GetFileReq"
}

func (m GetFileReq) PropertiesMap() map[string]any {
	return map[string]any{
		"filename": m.Filename,
		"ofs":      m.Offset, // TODO offset
		"len":      m.Length, // TODO length
	}
}
