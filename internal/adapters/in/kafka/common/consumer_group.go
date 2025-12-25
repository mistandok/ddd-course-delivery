package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"delivery/internal/pkg/errs"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

var ErrIncorrectMessage = errors.New("incorrect message format")

type EventHandler[TEvent proto.Message] interface {
	Handle(ctx context.Context, event TEvent) error
}

// KafkaConsumer объединяет consumer и handler в одной структуре
type KafkaConsumer[TEvent proto.Message] struct {
	topic         string
	consumerGroup sarama.ConsumerGroup
	domainHandler EventHandler[TEvent]
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewKafkaConsumerGroup[TEvent proto.Message](
	brokers []string,
	group string,
	topic string,
	domainHandler EventHandler[TEvent],
) (*KafkaConsumer[TEvent], error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if len(group) == 0 {
		return nil, errs.NewValueIsRequiredError("group")
	}
	if len(topic) == 0 {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V3_4_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer[TEvent]{
		topic:         topic,
		consumerGroup: consumerGroup,
		domainHandler: domainHandler,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (c *KafkaConsumer[TEvent]) Close() error {
	c.cancel()
	return c.consumerGroup.Close()
}

func (c *KafkaConsumer[TEvent]) Consume() error {
	for {
		err := c.consumerGroup.Consume(c.ctx, []string{c.topic}, c)
		if err != nil {
			return fmt.Errorf("failed to consume messages: %w", err)
		}
		if c.ctx.Err() != nil {
			return nil
		}
	}
}

func (c *KafkaConsumer[TEvent]) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *KafkaConsumer[TEvent]) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *KafkaConsumer[TEvent]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := context.Background()
		log.Printf("Received: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		eventType := reflect.TypeOf((*TEvent)(nil)).Elem()
		eventValue := reflect.New(eventType.Elem())
		event := eventValue.Interface().(TEvent)

		if err := proto.Unmarshal(message.Value, event); err != nil {
			return fmt.Errorf("%w: %v", ErrIncorrectMessage, err)
		}

		err := c.domainHandler.Handle(ctx, event)
		if err != nil {
			if errors.Is(err, ErrIncorrectMessage) {
				session.MarkMessage(message, "")
				continue
			}

			continue
		}

		// Mark message as consumed
		session.MarkMessage(message, "")
	}
	return nil
}
