package di

import (
	"go-payments-api/internal/infrastructure/messaging/kafka"
	"go-payments-api/internal/settings"

	"github.com/google/wire"
)

var messagingSet = wire.NewSet(
	provideKafkaPublisher,
)

func provideKafkaPublisher() kafka.Publisher {
	return kafka.NewPublisher(settings.Settings.Kafka.Brokers)
}
