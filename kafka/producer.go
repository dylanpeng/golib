package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/dylanpeng/golib/logger"
	"sync"
)

type ProducerConfig struct {
	Brokers []string `toml:"brokers" json:"brokers"`
}

type Producer struct {
	c      *ProducerConfig
	client sarama.AsyncProducer
	logger *logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func (p *Producer) Stop() {
	p.cancel()

	if err := p.client.Close(); err != nil {
		p.logger.Errorf("kafka producer close failed | brokers: %+v | error: %s", p.c.Brokers, err)
	}

	p.wg.Wait()
}

func (p *Producer) Send(topic string, body any, key string) error {
	payload, err := json.Marshal(body)

	if err != nil {
		p.logger.Errorf("send msg failed. | body: %+v | err: %s", body, err)
		return err
	}

	select {
	case <-p.ctx.Done():
		return errors.New("producer is stopped")
	default:
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(payload),
		}

		if key != "" {
			msg.Key = sarama.StringEncoder(key)
		}

		p.client.Input() <- msg
		p.logger.Debugf("send message success. | topic: %s | value: %+v", topic, body)
	}

	return nil
}

func (p *Producer) logErr() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case err := <-p.client.Errors():
			p.logger.Errorf("kafka producer receive error | brokers: %+v | error: %s", p.c.Brokers, err)
		}
	}
}

func NewProducer(c *ProducerConfig, logger *logger.Logger) (producer *Producer, err error) {
	producer = &Producer{
		c:      c,
		logger: logger,
	}

	producer.ctx, producer.cancel = context.WithCancel(context.Background())
	producer.wg = &sync.WaitGroup{}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = false // 设定是否需要返回成功信息
	config.Producer.Return.Errors = true     // 设定是否需要返回错误信息
	// ack等级
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区选择器
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer.client, err = sarama.NewAsyncProducer(c.Brokers, config)

	if err != nil {
		return
	}

	producer.wg.Add(1)
	go producer.logErr()

	return
}
