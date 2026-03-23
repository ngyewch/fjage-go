package gateway

import (
	"context"
	"io"

	"github.com/F2077/go-pubsub/pubsub"
)

type Transport interface {
	io.Closer

	SubscribeToRequests() (JsonMessageSubscription, error)
	SubscribeToResponse(jsonMessage *JSONMessage) (JsonMessageSubscription, error)
	SubscribeToMessageResponse(msgID string) (JsonMessageSubscription, error)
	SendJsonMessage(ctx context.Context, jsonMessage *JSONMessage) error
}

type JsonMessageSubscription interface {
	io.Closer

	Chan() <-chan *JSONMessage
	ErrChan() <-chan error
}

type PubSubJsonMessageSubscription struct {
	subscription *pubsub.Subscription[*JSONMessage]
}

func NewPubSubJsonMessageSubscription(subscription *pubsub.Subscription[*JSONMessage]) *PubSubJsonMessageSubscription {
	return &PubSubJsonMessageSubscription{
		subscription: subscription,
	}
}

func (subscription *PubSubJsonMessageSubscription) Close() error {
	// HACK sometimes this hangs
	go func() {
		_ = subscription.subscription.Close()
	}()
	return nil
}

func (subscription *PubSubJsonMessageSubscription) Chan() <-chan *JSONMessage {
	return subscription.subscription.Ch
}

func (subscription *PubSubJsonMessageSubscription) ErrChan() <-chan error {
	return subscription.subscription.ErrCh
}
