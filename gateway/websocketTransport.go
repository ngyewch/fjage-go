package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/F2077/go-pubsub/pubsub"
	"github.com/coder/websocket"
)

type WebSocketTransport struct {
	conn       *websocket.Conn
	broker     *pubsub.Broker[*JSONMessage]
	publisher  *pubsub.Publisher[*JSONMessage]
	subscriber *pubsub.Subscriber[*JSONMessage]
	closed     bool
}

func NewWebSocketTransport(ctx context.Context, gatewayUrl string) (*WebSocketTransport, error) {
	u, err := url.Parse(gatewayUrl)
	if err != nil {
		return nil, err
	}
	if (u.Scheme != "ws") && (u.Scheme != "wss") {
		return nil, fmt.Errorf("invalid websocket gateway url")
	}
	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	conn, resp, err := websocket.Dial(dialCtx, gatewayUrl, nil)
	if (resp != nil) && (resp.Body != nil) {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	broker, err := pubsub.NewBroker[*JSONMessage]()
	if err != nil {
		return nil, err
	}
	transport := &WebSocketTransport{
		conn:       conn,
		broker:     broker,
		publisher:  pubsub.NewPublisher[*JSONMessage](broker),
		subscriber: pubsub.NewSubscriber[*JSONMessage](broker),
	}
	go transport.keepAlive()
	go transport.readLoop()
	return transport, nil
}

func (transport *WebSocketTransport) Close() error {
	_ = transport.subscriber.Close()
	_ = transport.conn.CloseNow()
	transport.closed = true
	return nil
}

func (transport *WebSocketTransport) SubscribeToRequests() (JsonMessageSubscription, error) {
	subscription, err := transport.subscriber.Subscribe("requests",
		pubsub.WithChannelSize[*JSONMessage](pubsub.DefaultChannelSize),
	)
	if err != nil {
		return nil, err
	}
	return NewPubSubJsonMessageSubscription(subscription), nil
}

func (transport *WebSocketTransport) SubscribeToResponse(req *JSONMessage) (JsonMessageSubscription, error) {
	subscription, err := transport.subscriber.Subscribe(fmt.Sprintf("response/%s/%s", req.Action, req.ID),
		pubsub.WithChannelSize[*JSONMessage](pubsub.Single),
	)
	if err != nil {
		return nil, err
	}
	return NewPubSubJsonMessageSubscription(subscription), nil
}

func (transport *WebSocketTransport) SubscribeToMessageResponse(msgID string) (JsonMessageSubscription, error) {
	topic := fmt.Sprintf("messageResponse/%s", msgID)
	subscription, err := transport.subscriber.Subscribe(topic,
		pubsub.WithChannelSize[*JSONMessage](pubsub.Single),
	)
	if err != nil {
		return nil, err
	}
	return NewPubSubJsonMessageSubscription(subscription), nil
}

func (transport *WebSocketTransport) SendJsonMessage(ctx context.Context, jsonMessage *JSONMessage) error {
	if jsonMessage.ID == "" {
		return fmt.Errorf("JSONMessage.id is empty")
	}
	jsonMessageBytes, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}
	jsonMessageBytes = append(jsonMessageBytes, '\n')
	//fmt.Println(time.Now().Format(time.DateTime), ">>>", string(jsonMessageBytes))
	return transport.conn.Write(ctx, websocket.MessageText, jsonMessageBytes)
}

func (transport *WebSocketTransport) readLoop() {
	ctx := context.Background()
	for !transport.closed {
		messageType, messageBytes, err := transport.conn.Read(ctx)
		if err != nil {
			/*
				closeStatus := websocket.CloseStatus(err)
				if closeStatus == -1 {
					slog.Error("websocket read error",
						slog.Any("err", err),
						slog.Any("closeStatus", closeStatus),
					)
				}
			*/
			break
		} else {
			if messageType != websocket.MessageText {
				continue
			}
			//fmt.Println(time.Now().Format(time.DateTime), "<<<", string(messageBytes))
			var jsonMessage JSONMessage
			err = json.Unmarshal(messageBytes, &jsonMessage)
			if err != nil {
				slog.Error("JSONMessage unmarshal error",
					slog.Any("err", err),
				)
				continue
			}
			topic := "requests"
			if jsonMessage.InResponseTo != "" {
				topic = fmt.Sprintf("response/%s/%s", jsonMessage.InResponseTo, jsonMessage.ID)
			} else if (jsonMessage.Action == "send") && (jsonMessage.Message != nil) && (jsonMessage.Message.Data != nil) {
				inReplyTo, ok := jsonMessage.Message.Data["inReplyTo"].(string)
				if ok {
					topic = fmt.Sprintf("messageResponse/%s", inReplyTo)
				}
			}
			err = transport.publisher.Publish(topic, &jsonMessage)
			if err != nil {
				slog.Error("JSONMessage publish error",
					slog.Any("err", err),
				)
			}
		}
	}
}

type AliveMessage struct {
	Alive bool `json:"alive"`
}

func (transport *WebSocketTransport) keepAlive() {
	aliveMessageBytes, err := json.Marshal(AliveMessage{Alive: true})
	if err != nil {
		slog.Error("alive message marshal error",
			slog.Any("err", err),
		)
		return
	}
	aliveMessageBytes = append(aliveMessageBytes, '\n')
	ctx := context.Background()
	for {
		err = transport.conn.Write(ctx, websocket.MessageText, aliveMessageBytes)
		if err != nil {
			slog.Warn("error sending alive message",
				slog.Any("err", err),
			)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}
