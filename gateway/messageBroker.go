package gateway

import (
	"github.com/F2077/go-pubsub/pubsub"
)

type MessageBroker[T any] struct {
	broker    *pubsub.Broker[T]
	publisher *pubsub.Publisher[T]
}

func NewMessageBroker[T any]() (*MessageBroker[T], error) {
	broker, err := pubsub.NewBroker[T]()
	if err != nil {
		return nil, err
	}
	publisher := pubsub.NewPublisher[T](broker)

	return &MessageBroker[T]{
		broker:    broker,
		publisher: publisher,
	}, nil
}

func (messageBroker *MessageBroker[T]) Close() error {
	return nil
}

func (messageBroker *MessageBroker[T]) Subscribe(topic string) (*MessageSubscription[T], error) {
	subscriber := pubsub.NewSubscriber[T](messageBroker.broker)
	subscription, err := subscriber.Subscribe(topic, pubsub.WithChannelSize[T](pubsub.DefaultChannelSize))
	if err != nil {
		_ = subscriber.Close()
		return nil, err
	}
	return &MessageSubscription[T]{
		Ch:           subscription.Ch,
		ErrCh:        subscription.ErrCh,
		subscriber:   subscriber,
		subscription: subscription,
	}, nil

}

func (messageBroker *MessageBroker[T]) Publish(topic string, msg T) error {
	return messageBroker.publisher.Publish(topic, msg)
}

type MessageSubscription[T any] struct {
	Ch           <-chan T
	ErrCh        <-chan error
	subscriber   *pubsub.Subscriber[T]
	subscription *pubsub.Subscription[T]
}

func (subscription *MessageSubscription[T]) Close() error {
	_ = subscription.subscription.Close()
	_ = subscription.subscriber.Close()
	return nil
}
