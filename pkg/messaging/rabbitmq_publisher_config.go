package messaging

import "strings"

type PublisherConfig interface {
	Topic() string
	Format() string
}

type RabbitMQPublisherConfigFunc func(*RabbitMQPublisherConfig)

type RabbitMQPublisherConfig struct {
	topic         string
	serviceName   string
	messageFormat string
}

func newDefaultRabbitMQPublisherConfig() *RabbitMQPublisherConfig {
	return &RabbitMQPublisherConfig{
		topic:         "random_topic",
		serviceName:   "random_service_name",
		messageFormat: "text/plain",
	}
}

// WithTopic sets the exchange name to publish the messages
// On RabbitMQ, the exchange is the entity that receives messages from the publisher
// and routes them to the queues
func WithTopic(topic string) RabbitMQPublisherConfigFunc {
	return func(config *RabbitMQPublisherConfig) {
		config.topic = topic
	}
}

// WithServiceName sets the service name that is publishing the messages
func WithServiceName(serviceName string) RabbitMQPublisherConfigFunc {
	return func(config *RabbitMQPublisherConfig) {
		config.serviceName = strings.ToLower(serviceName)
	}
}

// WithJSONMessageFormat sets the message format as application/json
func WithJSONMessageFormat() RabbitMQPublisherConfigFunc {
	return func(config *RabbitMQPublisherConfig) {
		config.messageFormat = "application/json"
	}
}

func NewRabbitMQPublisherConfig(options ...RabbitMQPublisherConfigFunc) *RabbitMQPublisherConfig {
	config := newDefaultRabbitMQPublisherConfig()
	for _, option := range options {
		option(config)
	}

	return config
}

func (rpc *RabbitMQPublisherConfig) Topic() string {
	return rpc.topic
}

func (rpc *RabbitMQPublisherConfig) ServiceName() string {
	return rpc.serviceName
}

func (rpc *RabbitMQPublisherConfig) Format() string {
	return rpc.messageFormat
}
