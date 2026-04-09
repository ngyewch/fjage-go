package gateway

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ngyewch/fjage-go"
)

type DefaultGateway struct {
	transport    Transport
	subscription JsonMessageSubscription
	agentID      string
}

func NewDefaultGateway(ctx context.Context, transport Transport) (*DefaultGateway, error) {
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
	wantMessagesForJSONMessage, err := NewWantsMessagesFor([]string{gw.agentID})
	if err != nil {
		return nil, err
	}
	err = gw.transport.SendJsonMessage(ctx, wantMessagesForJSONMessage)
	if err != nil {
		return nil, err
	}
	go gw.messageHandler()
	return gw, nil
}

func (gw *DefaultGateway) Close() error {
	_ = gw.subscription.Close()
	return nil
}

func (gw *DefaultGateway) Transport() Transport {
	return gw.transport
}

func (gw *DefaultGateway) AgentID() string {
	return gw.agentID
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
				break
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
				// do nothing
				/*
					slog.Debug("unhandled message",
						slog.Any("jsonMessage", jsonMessage),
					)
				*/
			}
		}
	}
}

func (gw *DefaultGateway) request(ctx context.Context, req *JSONMessage) (*JSONMessage, error) {
	subscription, err := gw.transport.SubscribeToResponse(req)
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
		case rsp := <-subscription.Chan():
			if (rsp.ID == req.ID) && (rsp.InResponseTo == req.Action) {
				return rsp, nil
			}
		}
	}
}

func (gw *DefaultGateway) requestSend(ctx context.Context, req *JSONMessage, msgID string) (*JSONMessage, error) {
	subscription, err := gw.transport.SubscribeToMessageResponse(msgID)
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
		case rsp := <-subscription.Chan():
			if (rsp.Message == nil) || (rsp.Message.Data == nil) {
				continue
			}
			inReplyTo, ok := rsp.Message.Data["inReplyTo"].(string)
			if ok && (inReplyTo == msgID) {
				return rsp, nil
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

func (gw *DefaultGateway) Send(ctx context.Context, clazz string, message *fjage.Message, properties map[string]any) (*SendResponse, error) {
	req, err := NewSendRequestMessage(clazz, message, properties)
	if err != nil {
		return nil, err
	}
	rsp, err := gw.requestSend(ctx, req, message.MsgID)
	if err != nil {
		return nil, err
	}
	return &SendResponse{
		Message: rsp.Message,
	}, nil
}
