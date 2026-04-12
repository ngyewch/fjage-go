package fjage

type IMessage interface {
	JavaClassName() string
	Header() *Message
}

type Message struct {
	MsgID        string       `json:"msgID,omitempty"`
	Performative Performative `json:"perf,omitempty"`
	Recipient    string       `json:"recipient,omitempty"`
	Sender       string       `json:"sender,omitempty"`
	InReplyTo    string       `json:"inReplyTo,omitempty"`
	SentAt       int64        `json:"sentAt,omitempty"`
}

func (m *Message) JavaClassName() string {
	return "org.arl.fjage.Message"
}

func (m *Message) Header() *Message {
	return m
}
