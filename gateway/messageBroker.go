package gateway

import (
	"github.com/F2077/go-pubsub/pubsub"
)

type MessageBroker[T any] struct {
	broker     *pubsub.Broker[T]
	publisher  *pubsub.Publisher[T]
	subscriber *pubsub.Subscriber[T]
}

func NewMessageBroker[T any]() (*MessageBroker[T], error) {
	broker, err := pubsub.NewBroker[T]()
	if err != nil {
		return nil, err
	}
	publisher := pubsub.NewPublisher[T](broker)
	subscriber := pubsub.NewSubscriber[T](broker)

	return &MessageBroker[T]{
		broker:     broker,
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

func (messageBroker *MessageBroker[T]) Close() error {
	_ = messageBroker.subscriber.Close()
	return nil
}

func (messageBroker *MessageBroker[T]) Subscribe(topic string) (*MessageSubscription[T], error) {
	subscription, err := messageBroker.subscriber.Subscribe(topic, pubsub.WithChannelSize[T](pubsub.DefaultChannelSize))
	if err != nil {
		return nil, err
	}
	return &MessageSubscription[T]{
		Ch:           subscription.Ch,
		ErrCh:        subscription.ErrCh,
		subscription: subscription,
	}, nil

}

func (messageBroker *MessageBroker[T]) Publish(topic string, msg T) error {
	return messageBroker.publisher.Publish(topic, msg)
}

type MessageSubscription[T any] struct {
	Ch           <-chan T
	ErrCh        <-chan error
	subscription *pubsub.Subscription[T]
}

func (subscription *MessageSubscription[T]) Close() error {
	return subscription.subscription.Close()
}
