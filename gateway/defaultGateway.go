package gateway

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type DefaultGateway struct {
	transport    Transport
	subscription JsonMessageSubscription
	agentID      string
}

func NewDefaultGateway(transport Transport) (*DefaultGateway, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	subscription, err := transport.SubscribeToRequests()
	if err != nil {
		return nil, err
	}
	gw := &DefaultGateway{
		transport:    transport,
		subscription: subscription,
		agentID:      "gateway-" + id.String(),
	}
	go gw.messageHandler()
	return gw, nil
}

func (gw *DefaultGateway) Close() error {
	_ = gw.subscription.Close()
	return nil
}

func (gw *DefaultGateway) messageHandler() {
	for {
		select {
		case err := <-gw.subscription.ErrChan():
			if err != nil {
				slog.Error("error reading from subscription",
					slog.Any("err", err),
				)
				break
			}
		case jsonMessage := <-gw.subscription.Chan():
			if jsonMessage == nil {
				break
			}
			if jsonMessage.InResponseTo != "" {
				continue
			}
			switch jsonMessage.Action {
			case "agents":
				rsp := NewAgentsResponseMessage(jsonMessage, []string{gw.agentID})
				err := gw.transport.SendJsonMessage(context.Background(), rsp)
				if err != nil {
					slog.Error("error sending response",
						slog.Any("err", err),
					)
					break
				}
			case "services":
				rsp := NewServicesResponseMessage(jsonMessage, make([]string, 0))
				err := gw.transport.SendJsonMessage(context.Background(), rsp)
				if err != nil {
					slog.Error("error sending response",
						slog.Any("err", err),
					)
					break
				}
			case "agentForService":
				rsp := NewAgentForServiceResponseMessage(jsonMessage, "")
				err := gw.transport.SendJsonMessage(context.Background(), rsp)
				if err != nil {
					slog.Error("error sending response",
						slog.Any("err", err),
					)
					break
				}
			case "agentsForService":
				rsp := NewAgentsForServiceResponseMessage(jsonMessage, make([]string, 0))
				err := gw.transport.SendJsonMessage(context.Background(), rsp)
				if err != nil {
					slog.Error("error sending response",
						slog.Any("err", err),
					)
					break
				}
			case "containsAgent":
				rsp := NewContainsAgentResponseMessage(jsonMessage, jsonMessage.AgentID == gw.agentID)
				err := gw.transport.SendJsonMessage(context.Background(), rsp)
				if err != nil {
					slog.Error("error sending response",
						slog.Any("err", err),
					)
					break
				}
			default:
				fmt.Printf("!!! unhandled request %+v\n", jsonMessage)
			}
		}
	}
}

func (gw *DefaultGateway) request(ctx context.Context, req *JSONMessage) (*JSONMessage, error) {
	subscription, err := gw.transport.SubscribeToResponse(req.ID, req.Action)
	if err != nil {
		return nil, err
	}
	defer func(subscription JsonMessageSubscription) {
		_ = subscription.Close()
	}(subscription)
	err = gw.transport.SendJsonMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-subscription.ErrChan():
			return nil, err
		case jsonMessage := <-subscription.Chan():
			if (jsonMessage.ID == req.ID) && (jsonMessage.InResponseTo == req.Action) {
				return jsonMessage, nil
			}
		}
	}
}

func (gw *DefaultGateway) Agents(ctx context.Context) (*AgentsResponse, error) {
	req, err := NewAgentsRequestMessage()
	if err != nil {
		return nil, err
	}
	rsp, err := gw.request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &AgentsResponse{
		AgentIDs: rsp.AgentIDs,
	}, nil
}

func (gw *DefaultGateway) Services(ctx context.Context) (*ServicesResponse, error) {
	req, err := NewServicesRequestMessage()
	if err != nil {
		return nil, err
	}
	rsp, err := gw.request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ServicesResponse{
		Services: rsp.Services,
	}, nil
}

func (gw *DefaultGateway) AgentForService(ctx context.Context, service string) (*AgentForServiceResponse, error) {
	req, err := NewAgentForServiceRequestMessage(service)
	if err != nil {
		return nil, err
	}
	rsp, err := gw.request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &AgentForServiceResponse{
		AgentID: rsp.AgentID,
	}, nil
}

func (gw *DefaultGateway) AgentsForService(ctx context.Context, service string) (*AgentsForServiceResponse, error) {
	req, err := NewAgentsForServiceRequestMessage(service)
	if err != nil {
		return nil, err
	}
	rsp, err := gw.request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &AgentsForServiceResponse{
		AgentIDs: rsp.AgentIDs,
	}, nil
}

func (gw *DefaultGateway) ContainsAgent(ctx context.Context, agentID string) (*ContainsAgentResponse, error) {
	req, err := NewContainsAgentRequestMessage(agentID)
	if err != nil {
		return nil, err
	}
	rsp, err := gw.request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ContainsAgentResponse{
		Answer: rsp.Answer,
	}, nil
}
