package rocketmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	oProducer "github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/dylanpeng/golib/logger"
)

type Message struct {
	Topic       string   `json:"topic" toml:"topic" yaml:"topic"`
	Payload     []byte   `json:"payload" toml:"payload" yaml:"payload"`
	Tag         string   `json:"tag" toml:"tag" yaml:"tag"`
	Keys        []string `json:"keys" toml:"keys" yaml:"keys"`
	ShardingKey string   `json:"sharding_key" toml:"sharding_key" yaml:"sharding_key"`
	DelayLevel  int      `json:"delay_level" toml:"delay_level" yaml:"delay_level"`
}

func (m *Message) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("%+v", *m)
}

func (m *Message) Request() *primitive.Message {
	msg := primitive.NewMessage(m.Topic, m.Payload)

	if m.Tag != "" {
		msg.WithTag(m.Tag)
	}

	if len(m.Keys) > 0 {
		msg.WithKeys(m.Keys)
	}

	if m.ShardingKey != "" {
		msg.WithShardingKey(m.ShardingKey)
	}

	if m.DelayLevel > 0 {
		msg.WithDelayTimeLevel(m.DelayLevel)
	}

	return msg
}

type ProducerConfig struct {
	*ConnectConfig
	Group      string `toml:"group" json:"group"`
	RetryTimes int    `toml:"retry_times" json:"retry_times"`
}

type Producer struct {
	c        *ProducerConfig
	producer rocketmq.Producer
	ctx      context.Context
	logger   logger.ILogger
	cancel   context.CancelFunc
}

func (c *Producer) closed() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

func (c *Producer) Stop() {
	c.cancel()

	if err := c.producer.Shutdown(); err != nil {
		c.logger.Errorf("stop rocketmq producer failed | error: %s", err)
	}
}

func (c *Producer) SendSync(msg *Message) (result *primitive.SendResult, err error) {
	if c.closed() {
		err = errors.New("producer is stopped")
		return
	}

	result, err = c.producer.SendSync(c.ctx, msg.Request())

	if err != nil {
		c.logger.Errorf("rocketmq SendSync fail. |  msg: %s | err: %s", msg, err)
	} else {
		c.logger.Debugf("rocketmq SendSync message | message: %+v | result: %s", msg, result)
	}

	return
}

func (c *Producer) SendAsync(msg *Message) (err error) {
	if c.closed() {
		err = errors.New("producer is stopped")
		return
	}

	err = c.producer.SendAsync(c.ctx, func(_ context.Context, res *primitive.SendResult, e error) {
		if e != nil {
			c.logger.Errorf("send rocketmq message failed | message: %s | result: %s | error: %s", msg, res, e)
		} else {
			c.logger.Debugf("send rocketmq message | message: %+v | result: %s", msg, res)
		}

	}, msg.Request())
	return
}

func NewProducer(c *ProducerConfig, logger logger.ILogger) (producer *Producer, err error) {
	opts := []oProducer.Option{
		oProducer.WithNsResolver(primitive.NewPassthroughResolver(c.Endpoints)),
	}

	if c.Group != "" {
		opts = append(opts, oProducer.WithGroupName(c.Group))
	}

	if c.AccessKey != "" && c.SecretKey != "" {
		opts = append(opts, oProducer.WithCredentials(primitive.Credentials{
			AccessKey:     c.AccessKey,
			SecretKey:     c.SecretKey,
			SecurityToken: c.SecurityToken,
		}))
	}

	if c.RetryTimes > 0 {
		opts = append(opts, oProducer.WithRetry(c.RetryTimes))
	}

	rlog.SetLogger(&Logger{
		logger: logger,
		quiet:  true,
	})

	p, err := rocketmq.NewProducer(opts...)

	if err != nil {
		return
	}

	if err = p.Start(); err != nil {
		return
	}

	producer = &Producer{c: c, logger: logger, producer: p}
	producer.ctx, producer.cancel = context.WithCancel(context.Background())
	return
}
