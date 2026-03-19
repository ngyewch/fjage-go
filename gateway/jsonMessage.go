package gateway

import "github.com/google/uuid"

type JSONMessage struct {
	ID           string   `json:"id,omitempty"`
	Action       string   `json:"action,omitempty"`
	InResponseTo string   `json:"inResponseTo,omitempty"`
	AgentID      string   `json:"agentID,omitempty"`
	AgentIDs     []string `json:"agentIDs,omitempty"`
	AgentTypes   []string `json:"agentTypes,omitempty"`
	Service      string   `json:"service,omitempty"`
	Services     []string `json:"services,omitempty"`
	Answer       bool     `json:"answer,omitempty"`
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
		ID:       req.ID,
		Action:   req.Action,
		AgentIDs: agentIDs,
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
		ID:       req.ID,
		Action:   req.Action,
		Services: services,
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
		ID:      req.ID,
		Action:  req.Action,
		AgentID: agentID,
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
		ID:       req.ID,
		Action:   req.Action,
		AgentIDs: agentIDs,
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
		ID:     req.ID,
		Action: req.Action,
		Answer: answer,
	}
}
