package env

import (
	"delivery/internal/config"
	"errors"
	"os"
)

const (
	kafkaHost                 = "KAFKA_HOST"
	kafkaConsumerGroup        = "KAFKA_CONSUMER_GROUP"
	kafkaBasketConfirmedTopic = "KAFKA_BASKET_CONFIRMED_TOPIC"
	kafkaOrderChangedTopic    = "KAFKA_ORDER_CHANGED_TOPIC"
)

type KafkaCfgSearcher struct{}

func NewKafkaCfgSearcher() *KafkaCfgSearcher {
	return &KafkaCfgSearcher{}
}

func (k *KafkaCfgSearcher) Get() (*config.KafkaConfig, error) {
	host := os.Getenv(kafkaHost)
	if len(host) == 0 {
		return nil, errors.New("KAFKA_HOST is not set")
	}

	consumerGroup := os.Getenv(kafkaConsumerGroup)
	if len(consumerGroup) == 0 {
		return nil, errors.New("KAFKA_CONSUMER_GROUP is not set")
	}

	basketConfirmedTopic := os.Getenv(kafkaBasketConfirmedTopic)
	if len(basketConfirmedTopic) == 0 {
		return nil, errors.New("KAFKA_BASKET_CONFIRMED_TOPIC is not set")
	}

	orderChangedTopic := os.Getenv(kafkaOrderChangedTopic)
	if len(orderChangedTopic) == 0 {
		return nil, errors.New("KAFKA_ORDER_CHANGED_TOPIC is not set")
	}

	return &config.KafkaConfig{
		Host:                 host,
		ConsumerGroup:        consumerGroup,
		BasketConfirmedTopic: basketConfirmedTopic,
		OrderChangedTopic:    orderChangedTopic,
	}, nil
}
