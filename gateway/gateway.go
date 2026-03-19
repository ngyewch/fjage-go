package gateway

import (
	"context"
	"io"

	"github.com/ngyewch/fjage-go"
)

type Gateway interface {
	io.Closer

	AgentID() string
	Agents(ctx context.Context) (*AgentsResponse, error)
	Services(ctx context.Context) (*ServicesResponse, error)
	AgentForService(ctx context.Context, service string) (*AgentForServiceResponse, error)
	AgentsForService(ctx context.Context, service string) (*AgentsForServiceResponse, error)
	ContainsAgent(ctx context.Context, agentID string) (*ContainsAgentResponse, error)
	Send(ctx context.Context, clazz string, message *fjage.Message, properties map[string]any) (*SendResponse, error)
}

type AgentsResponse struct {
	AgentIDs []string
}

type ServicesResponse struct {
	Services []string
}

type AgentForServiceResponse struct {
	AgentID string
}

type AgentsForServiceResponse struct {
	AgentIDs []string
}

type ContainsAgentResponse struct {
	Answer bool
}

type SendResponse struct {
	Message *MessageEnvelope
}
