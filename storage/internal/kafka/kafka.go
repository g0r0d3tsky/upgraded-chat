package kafka

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/config"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/IBM/sarama"
)

type MessageHandler func(*sarama.ConsumerMessage) error

type Consumer struct {
	handler MessageHandler
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	//
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	//
	return nil
}

func Init() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	return cfg
}

func NewConsumer(handler MessageHandler) *Consumer {
	return &Consumer{handler: handler}
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.Info("message channel closed")
				return nil
			}
			err := c.handler(message)
			if err != nil {
				slog.Error("handle message", err)
				return fmt.Errorf("handle message: %w", err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			session.Commit()
			return nil
		}
	}
}

func RunConsumer(ctx context.Context, wg *sync.WaitGroup, config *config.Config, consumer *Consumer) (sarama.ConsumerGroup, error) {
	consumerGroup, err := sarama.NewConsumerGroup(
		strings.Split(config.Kafka.BrokerList, ","), config.Kafka.GroupID, Init())
	if err != nil {
		slog.Error("creating consumer group:", err)
		return nil, err
	}
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, strings.Split(config.Kafka.Topics, ","), consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				slog.Error("from consumer:", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return consumerGroup, nil
}
