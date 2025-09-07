package di

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/soulcodex/deus-cargo-tracker/configs"
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/messaging"
)

func initEventPublisher(conn *amqp.Connection, cfg *configs.Config) domain.EventPublisher {
	publisherConfig := messaging.NewRabbitMQPublisherConfig(
		messaging.WithServiceName(cfg.AppServiceName),
		messaging.WithTopic(cfg.RabbitMQTopic),
		messaging.WithJSONMessageFormat(),
	)

	return messaging.NewRabbitMQPublisher(
		conn,
		publisherConfig,
	)
}

func initRabbitMQConnection(cfg *configs.Config) *amqp.Connection {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		panic(err)
	}

	return conn
}
