package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/dylanpeng/golib/logger"
	"sync"
)

type BaseConsumerGroupHandler struct {
	logger  *logger.Logger
	handler func([]byte) error
	ctx     context.Context
	worker  int
	wg      *sync.WaitGroup
}

func NewBaseConsumerGroupHandler(ctx context.Context, worker int, logger *logger.Logger, handler func([]byte) error) *BaseConsumerGroupHandler {
	result := &BaseConsumerGroupHandler{
		logger:  logger,
		handler: handler,
		ctx:     ctx,
		worker:  worker,
		wg:      &sync.WaitGroup{},
	}

	if worker <= 0 {
		result.worker = 1
	}

	return result
}

func (*BaseConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*BaseConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim every partition call once
func (c *BaseConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for i := 0; i < c.worker; i++ {
		c.wg.Add(1)
		go c.receive(sess, claim)
	}

	c.wg.Wait()
	c.logger.Infof("BaseConsumerGroupHandler ConsumeClaim finish.")
	return nil
}

func (c *BaseConsumerGroupHandler) receive(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("BaseConsumerGroupHandler receive context have close")
			return
		case msg, ok := <-claim.Messages():
			if !ok {
				c.logger.Infof("BaseConsumerGroupHandler receive message channel have close")
				return
			}

			if err := c.handler(msg.Value); err != nil {
				c.logger.Errorf("consume message failed. | Message topic:%q partition:%d offset:%d msg: %s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
			} else {
				c.logger.Infof("consume message. | Message topic:%q partition:%d offset:%d msg: %s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
				sess.MarkMessage(msg, "")
			}
		}
	}
}
