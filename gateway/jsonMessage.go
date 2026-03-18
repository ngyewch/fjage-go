package gateway

type JSONMessage struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

func NewAgentsMessage() *JSONMessage {
	return &JSONMessage{
		Action: "agents",
	}
}
