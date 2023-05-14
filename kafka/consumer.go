package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/dylanpeng/golib/logger"
)

type ConsumerConfig struct {
	Brokers []string `toml:"brokers" json:"brokers"`
	GroupId string   `toml:"group_id" json:"group_id"`
	Topic   string   `toml:"topic" json:"topic"`
	Worker  int      `toml:"worker" json:"worker"`
}

type Consumer struct {
	c       *ConsumerConfig
	client  sarama.ConsumerGroup
	logger  logger.ILogger
	ctx     context.Context
	handler sarama.ConsumerGroupHandler
}

func NewConsumer(c *ConsumerConfig, logger logger.ILogger, handle func([]byte) error) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始

	group, err := sarama.NewConsumerGroup(c.Brokers, c.GroupId, config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	consumer := &Consumer{
		c:       c,
		client:  group,
		logger:  logger,
		ctx:     ctx,
		handler: NewConsumerHandler(ctx, c.Worker, logger, handle),
	}

	go consumer.run()
	return consumer, nil
}

func (c *Consumer) run() {
	go c.logErr()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("Consumer run context close.")
			return
		default:
		}

		topics := []string{c.c.Topic}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := c.client.Consume(c.ctx, topics, c.handler)
		if err != nil {
			c.logger.Errorf("Consumer Consume Error: %s", err)
			return
		}
	}
}

// Track errors
func (c *Consumer) logErr() {
	for err := range c.client.Errors() {
		c.logger.Errorf("Error: %s", err)
	}
}

func (c *Consumer) Stop() {
	if err := c.client.Close(); err != nil {
		c.logger.Errorf("Consumer Stop Error: %s", err)
	}
}
