package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	oConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/dylanpeng/golib/logger"
	"strings"
	"sync"
)

type ConsumerConfig struct {
	*ConnectConfig
	Topic     string   `toml:"topic" json:"topic" yaml:"topic"`
	Group     string   `toml:"group" json:"group" yaml:"group"`
	Orderly   bool     `toml:"orderly" json:"orderly" yaml:"orderly"`
	FromFirst bool     `toml:"from_first" json:"from_first" yaml:"from_first"`
	Tags      []string `toml:"tags" json:"tags" yaml:"tags"`
	Worker    int      `toml:"worker" json:"worker" yaml:"worker"`
}

type Consumer struct {
	c        *ConsumerConfig
	handler  func([]byte) error
	consumer rocketmq.PushConsumer
	logger   logger.ILogger
	wg       *sync.WaitGroup
	msgQueue chan []byte
	mutex    sync.RWMutex
	close    bool
}

func (c *Consumer) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.consumer.Shutdown(); err != nil {
		c.logger.Errorf("rocketmq consumer shutdown failed | error: %s", err)
		return
	}

	c.close = true
	close(c.msgQueue)

	c.wg.Wait()
	c.logger.Debugf("stop consumer")
}

func (c *Consumer) run() {
	worker := c.c.Worker

	if worker <= 0 {
		worker = 1
	}

	for i := 0; i < worker; i++ {
		c.wg.Add(1)
		go c.startWork()
	}
}

func (c *Consumer) startWork() {
	defer c.wg.Done()

	for body := range c.msgQueue {
		if err := c.handler(body); err != nil {
			c.logger.Errorf("handle rocketmq message failed | message: %s | error: %s", string(body), err)
		} else {
			c.logger.Debugf("handle rocketmq message | message: %s", string(body))
		}
	}

	c.logger.Debugf("worker index: %d stop")
}

func (c *Consumer) receive(_ context.Context, msgs ...*primitive.MessageExt) (res oConsumer.ConsumeResult, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.close {
		c.logger.Debugf("consumer has closed")
		return oConsumer.ConsumeRetryLater, nil
	}

	for _, msg := range msgs {
		c.logger.Debugf("receive rocketmq message | message: %+v", msg)
		c.msgQueue <- msg.Body
	}

	return oConsumer.ConsumeSuccess, nil
}

func NewConsumer(conf *ConsumerConfig, handler func([]byte) error, logger *logger.Logger) (consumer *Consumer, err error) {
	opts := []oConsumer.Option{
		oConsumer.WithNsResolver(primitive.NewPassthroughResolver(conf.Endpoints)),
		oConsumer.WithConsumerModel(oConsumer.Clustering),
		oConsumer.WithGroupName(conf.Group),
	}

	if conf.Orderly {
		opts = append(opts, oConsumer.WithConsumerOrder(true))
	}

	if conf.FromFirst {
		opts = append(opts, oConsumer.WithConsumeFromWhere(oConsumer.ConsumeFromFirstOffset))
	} else {
		opts = append(opts, oConsumer.WithConsumeFromWhere(oConsumer.ConsumeFromLastOffset))
	}

	if conf.AccessKey != "" && conf.SecretKey != "" {
		opts = append(opts, oConsumer.WithCredentials(primitive.Credentials{
			AccessKey:     conf.AccessKey,
			SecretKey:     conf.SecretKey,
			SecurityToken: conf.SecurityToken,
		}))
	}

	selector := oConsumer.MessageSelector{}

	if len(conf.Tags) > 0 {
		selector = oConsumer.MessageSelector{
			Type:       oConsumer.TAG,
			Expression: strings.Join(conf.Tags, " || "),
		}
	}

	c := &Consumer{
		c:        conf,
		handler:  handler,
		logger:   logger,
		wg:       &sync.WaitGroup{},
		msgQueue: make(chan []byte),
	}

	rlog.SetLogger(&Logger{logger: logger, quiet: true})

	if c.consumer, err = rocketmq.NewPushConsumer(opts...); err != nil {
		return
	}

	if err = c.consumer.Subscribe(conf.Topic, selector, c.receive); err != nil {
		return
	}

	c.run()

	if err = c.consumer.Start(); err != nil {
		return
	}

	consumer = c
	return
}
