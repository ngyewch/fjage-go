package gateway

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ngyewch/fjage-go"
)

type JSONMessage struct {
	ID           string           `json:"id,omitempty"`
	Action       string           `json:"action,omitempty"`
	InResponseTo string           `json:"inResponseTo,omitempty"`
	AgentID      string           `json:"agentID,omitempty"`
	AgentIDs     []string         `json:"agentIDs,omitempty"`
	AgentTypes   []string         `json:"agentTypes,omitempty"`
	Service      string           `json:"service,omitempty"`
	Services     []string         `json:"services,omitempty"`
	Answer       bool             `json:"answer,omitempty"`
	Relay        bool             `json:"relay,omitempty"`
	Message      *MessageEnvelope `json:"message,omitempty"`
}

func (jm JSONMessage) LogValue() slog.Value {
	jsonBytes, err := json.Marshal(jm)
	jsonString := ""
	if err == nil {
		jsonString = string(jsonBytes)
	}
	return slog.StringValue(jsonString)
}

type MessageEnvelope struct {
	Clazz string         `json:"clazz,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

func NewAgentsRequestMessage() (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:     id.String(),
		Action: "agents",
	}, nil
}

func NewAgentsResponseMessage(req *JSONMessage, agentIDs []string) *JSONMessage {
	return &JSONMessage{
		ID:           req.ID,
		InResponseTo: req.Action,
		AgentIDs:     agentIDs,
	}
}

func NewServicesRequestMessage() (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:     id.String(),
		Action: "services",
	}, nil
}

func NewServicesResponseMessage(req *JSONMessage, services []string) *JSONMessage {
	return &JSONMessage{
		ID:           req.ID,
		InResponseTo: req.Action,
		Services:     services,
	}
}

func NewAgentForServiceRequestMessage(service string) (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:      id.String(),
		Action:  "agentForService",
		Service: service,
	}, nil
}

func NewAgentForServiceResponseMessage(req *JSONMessage, agentID string) *JSONMessage {
	return &JSONMessage{
		ID:           req.ID,
		InResponseTo: req.Action,
		AgentID:      agentID,
	}
}

func NewAgentsForServiceRequestMessage(service string) (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:      id.String(),
		Action:  "agentsForService",
		Service: service,
	}, nil
}

func NewAgentsForServiceResponseMessage(req *JSONMessage, agentIDs []string) *JSONMessage {
	return &JSONMessage{
		ID:           req.ID,
		InResponseTo: req.Action,
		AgentIDs:     agentIDs,
	}
}

func NewContainsAgentRequestMessage(agentID string) (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:      id.String(),
		Action:  "containsAgent",
		AgentID: agentID,
	}, nil
}

func NewContainsAgentResponseMessage(req *JSONMessage, answer bool) *JSONMessage {
	return &JSONMessage{
		ID:           req.ID,
		InResponseTo: req.Action,
		Answer:       answer,
	}
}

func NewWantsMessagesFor(agentIDs []string) (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &JSONMessage{
		ID:       id.String(),
		Action:   "wantsMessagesFor",
		AgentIDs: agentIDs,
	}, nil
}

func NewSendRequestMessage(clazz string, message *fjage.Message, properties map[string]any) (*JSONMessage, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	messageMap := make(map[string]any)
	for key, value := range properties {
		messageMap[key] = value
	}
	message.PopulateMap(messageMap)
	return &JSONMessage{
		ID:     id.String(),
		Action: "send",
		Message: &MessageEnvelope{
			Clazz: clazz,
			Data:  messageMap,
		},
		Relay: true,
	}, nil
}
