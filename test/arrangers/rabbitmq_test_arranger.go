package testarrangers

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQArranger struct {
	topic      string
	connection *amqp.Connection
}

func NewRabbitMQArranger(topic string, connection *amqp.Connection) *RabbitMQArranger {
	return &RabbitMQArranger{topic: topic, connection: connection}
}

func (rma *RabbitMQArranger) MustArrange(_ context.Context) {
	channel, channelErr := rma.connection.Channel()
	if nil != channelErr {
		panic(channelErr)
	}

	defer func() {
		_ = channel.Close()
	}()

	topicErr := channel.ExchangeDeclare(
		rma.topic,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)

	if topicErr != nil {
		panic(topicErr)
	}
}
