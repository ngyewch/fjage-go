package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/F2077/go-pubsub/pubsub"
	"github.com/coder/websocket"
)

const (
	DefaultWebSocketPingInterval = 5 * time.Second
	DefaultWebSocketPingTimeout  = 1 * time.Second
)

type WebSocketTransport struct {
	url        string
	options    WebSocketTransportOptions
	conn       *websocket.Conn
	broker     *pubsub.Broker[*JSONMessage]
	publisher  *pubsub.Publisher[*JSONMessage]
	subscriber *pubsub.Subscriber[*JSONMessage]
	closed     bool
}

type WebSocketTransportOptions struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
}

func NewWebSocketTransport(ctx context.Context, gatewayUrl string, options WebSocketTransportOptions) (*WebSocketTransport, error) {
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
		url:        gatewayUrl,
		options:    options,
		conn:       conn,
		broker:     broker,
		publisher:  pubsub.NewPublisher[*JSONMessage](broker),
		subscriber: pubsub.NewSubscriber[*JSONMessage](broker),
	}

	err = transport.sendAlive(ctx, true)
	if err != nil {
		return nil, err
	}

	go transport.keepAlive()

	go transport.readLoop()

	return transport, nil
}

func (transport *WebSocketTransport) Close() error {
	transport.closed = true
	_ = transport.subscriber.Close()
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = transport.sendAlive(ctxWithTimeout, false)
	_ = transport.conn.CloseNow()
	return nil
}

func (transport *WebSocketTransport) Url() string {
	return transport.url
}

func (transport *WebSocketTransport) IsClosed() bool {
	return transport.closed
}

func (transport *WebSocketTransport) SubscribeToRequests() (JsonMessageSubscription, error) {
	if transport.closed {
		return nil, ErrTransportClosed
	}
	subscription, err := transport.subscriber.Subscribe("requests",
		pubsub.WithChannelSize[*JSONMessage](pubsub.DefaultChannelSize),
	)
	if err != nil {
		return nil, err
	}
	return NewPubSubJsonMessageSubscription(subscription), nil
}

func (transport *WebSocketTransport) SubscribeToResponse(req *JSONMessage) (JsonMessageSubscription, error) {
	if transport.closed {
		return nil, ErrTransportClosed
	}
	subscription, err := transport.subscriber.Subscribe(fmt.Sprintf("response/%s/%s", req.Action, req.ID),
		pubsub.WithChannelSize[*JSONMessage](pubsub.Single),
	)
	if err != nil {
		return nil, err
	}
	return NewPubSubJsonMessageSubscription(subscription), nil
}

func (transport *WebSocketTransport) SubscribeToMessageResponse(msgID string) (JsonMessageSubscription, error) {
	if transport.closed {
		return nil, ErrTransportClosed
	}
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
	if transport.closed {
		return ErrTransportClosed
	}
	if jsonMessage.ID == "" {
		return fmt.Errorf("JSONMessage.id is empty")
	}
	jsonMessageBytes, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}
	jsonMessageBytes = append(jsonMessageBytes, '\n')
	return transport.send(ctx, jsonMessageBytes)
}

func (transport *WebSocketTransport) readLoop() {
	defer func() {
		slog.Debug("websocket readLoop exit")
	}()
	ctx := context.Background()
	for !transport.closed {
		messageType, messageBytes, err := transport.conn.Read(ctx)
		if err != nil {
			if !transport.closed {
				slog.Error("websocket read error",
					slog.Any("err", err),
				)
				_ = transport.Close()
			}
			break
		} else {
			if messageType != websocket.MessageText {
				continue
			}
			slog.Debug("wsTransport",
				slog.String("<<<", string(messageBytes)),
			)
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
				recipient, hasRecipient := jsonMessage.Message.Data["recipient"].(string)
				_, hasSender := jsonMessage.Message.Data["sender"].(string)
				if hasRecipient && hasSender && strings.HasPrefix(recipient, "#") {
					// notification message: use "requests" topic
				} else {
					inReplyTo, hasInReplyTo := jsonMessage.Message.Data["inReplyTo"].(string)
					if hasInReplyTo {
						topic = fmt.Sprintf("messageResponse/%s", inReplyTo)
					}
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

func (transport *WebSocketTransport) sendAlive(ctx context.Context, alive bool) error {
	aliveMessageBytes, err := json.Marshal(AliveMessage{Alive: alive})
	if err != nil {
		return err
	}
	aliveMessageBytes = append(aliveMessageBytes, '\n')
	return transport.send(ctx, aliveMessageBytes)
}

func (transport *WebSocketTransport) pingWithTimeout(ctx context.Context, timeout time.Duration) error {
	if timeout > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return transport.conn.Ping(ctxWithTimeout)
	}
	return transport.conn.Ping(ctx)
}

func (transport *WebSocketTransport) keepAlive() {
	defer func() {
		slog.Debug("websocket keepAlive exit")
	}()
	ctx := context.Background()
	for !transport.closed {
		err := transport.pingWithTimeout(ctx, transport.options.PingTimeout)
		if err != nil {
			if !transport.closed {
				slog.Warn("error sending ping",
					slog.Any("err", err),
				)
				_ = transport.Close()
			}
			break
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(transport.options.PingInterval):
		}
	}
}

func (transport *WebSocketTransport) send(ctx context.Context, data []byte) error {
	slog.Debug("wsTransport",
		slog.String(">>>", string(data)),
	)
	return transport.conn.Write(ctx, websocket.MessageText, data)
}
