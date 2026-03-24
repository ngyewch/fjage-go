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
	if m.MsgID != "" {
		properties["msgID"] = m.MsgID
	}
	if m.Perf != "" {
		properties["perf"] = m.Perf
	}
	if m.Recipient != "" {
		properties["recipient"] = m.Recipient
	}
	if m.Sender != "" {
		properties["sender"] = m.Sender
	}
	if m.InReplyTo != "" {
		properties["inReplyTo"] = m.InReplyTo
	}
	if m.SentAt != 0 {
		properties["sentAt"] = m.SentAt
	}
}
