package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/dylanpeng/golib/logger"
	"sync"
)

type ConsumerHandler struct {
	logger  logger.ILogger
	handler func([]byte) error
	ctx     context.Context
	worker  int
	wg      *sync.WaitGroup
}

func NewConsumerHandler(ctx context.Context, worker int, logger logger.ILogger, handler func([]byte) error) *ConsumerHandler {
	result := &ConsumerHandler{
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

func (*ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim every partition call once
func (c *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for i := 0; i < c.worker; i++ {
		c.wg.Add(1)
		go c.receive(sess, claim)
	}

	c.wg.Wait()
	c.logger.Infof("ConsumerHandler ConsumeClaim finish.")
	return nil
}

func (c *ConsumerHandler) receive(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("ConsumerHandler receive context have close")
			return
		case msg, ok := <-claim.Messages():
			if !ok {
				c.logger.Infof("ConsumerHandler receive message channel have close")
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
