package fjage

type Message struct {
	MsgID     string `json:"msgID,omitempty"`
	Perf      string `json:"perf,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Sender    string `json:"sender,omitempty"`
	InReplyTo string `json:"inReplyTo,omitempty"`
	SentAt    int64  `json:"sentAt,omitempty"`
}

func (m Message) PopulateMap(properties map[string]any) {
	properties["msgID"] = m.MsgID
	properties["perf"] = m.Perf
	properties["recipient"] = m.Recipient
	properties["sender"] = m.Sender
	properties["inReplyTo"] = m.InReplyTo
	properties["sentAt"] = m.SentAt
}
