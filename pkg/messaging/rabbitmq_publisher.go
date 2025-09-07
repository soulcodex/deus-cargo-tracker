package messaging

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrFailedToPublishMessage = errutil.NewError("rabbitmq publisher failed to publish message")
)

//go:generate moq -pkg messagingmock -out mock/messaging_publisher_moq.go . Publisher
type Publisher interface {
	Publish(ctx context.Context, messages ...Message) error
}
type RabbitMQPublisher struct {
	conn   *amqp.Connection
	config *RabbitMQPublisherConfig
}

func NewRabbitMQPublisher(conn *amqp.Connection, config *RabbitMQPublisherConfig) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		conn:   conn,
		config: config,
	}
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, messages ...Message) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return ErrFailedToPublishMessage.Wrap(err)
	}

	defer func() {
		_ = channel.Close()
	}()

	for _, msg := range messages {
		marshaledMsg, marshalErr := p.amqpMessageFromEvent(msg)
		if marshalErr != nil {
			return ErrFailedToPublishMessage.Wrap(marshalErr)
		}

		marshalErr = channel.PublishWithContext(
			ctx,
			p.config.Topic(),
			msg.Type(),
			false,
			false,
			marshaledMsg,
		)

		if marshalErr != nil {
			return ErrFailedToPublishMessage.Wrap(marshalErr)
		}
	}

	return nil
}

func (p *RabbitMQPublisher) amqpMessageFromEvent(msg Message) (amqp.Publishing, error) {
	serializedEvent, marshalErr := json.Marshal(msg)
	if marshalErr != nil {
		return amqp.Publishing{}, ErrFailedToPublishMessage.Wrap(marshalErr)
	}

	publishable := amqp.Publishing{
		Headers: amqp.Table{
			"service":     p.config.ServiceName(),
			"occurred_on": msg.Time().Format(time.RFC3339),
		},
		MessageId:   msg.Identifier(),
		ContentType: p.config.Format(),
		Body:        serializedEvent,
		Timestamp:   msg.Time(),
		AppId:       p.config.ServiceName(),
	}

	return publishable, nil
}
