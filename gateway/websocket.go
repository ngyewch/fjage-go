package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WebSocketGateway struct {
	conn   *websocket.Conn
	closed bool
}

func NewWebSocketGateway(ctx context.Context, gatewayUrl string) (*WebSocketGateway, error) {
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
	gw := &WebSocketGateway{
		conn: conn,
	}
	go gw.readLoop()
	return gw, nil
}

func (gw *WebSocketGateway) Close() error {
	_ = gw.conn.CloseNow()
	gw.closed = true
	return nil
}

func (gw *WebSocketGateway) readLoop() {
	ctx := context.Background()
	for !gw.closed {
		messageType, messageBytes, err := gw.conn.Read(ctx)
		if err != nil {
			closeStatus := websocket.CloseStatus(err)
			if closeStatus == -1 {
				slog.Error("websocket read error",
					slog.Any("err", err),
					slog.Any("closeStatus", closeStatus),
				)
			} else {
				break
			}
		} else {
			fmt.Printf("%s / %s\n", messageType.String(), string(messageBytes))
		}
	}
}

func (gw *WebSocketGateway) sendJsonMessage(ctx context.Context, jsonMessage *JSONMessage) error {
	if jsonMessage.ID == "" {
		jsonMessage.ID = uuid.New().String()
	}
	jsonMessageBytes, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}
	jsonMessageBytes = append(jsonMessageBytes, '\n')
	fmt.Println(">>>", string(jsonMessageBytes))
	return gw.conn.Write(ctx, websocket.MessageText, jsonMessageBytes)
}

func (gw *WebSocketGateway) Agents(ctx context.Context) error {
	return gw.sendJsonMessage(ctx, NewAgentsMessage())
}
