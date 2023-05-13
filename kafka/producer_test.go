package kafka

import (
	"github.com/dylanpeng/golib/logger"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestProducerConsumerGroup(t *testing.T) {
	producerConf := &ProducerConfig{
		Brokers: []string{"localhost:9092"},
	}

	Log, err := logger.NewLogger(&logger.Config{
		FilePath:   "./logs/confuse",
		Level:      "debug",
		TimeFormat: "2006-01-02 15:04:05.000",
		MaxAgeDay:  30,
	})

	if err != nil {
		t.Fatalf("NewLogger fail. | err: %s", err)
	}

	producer, err := NewProducer(producerConf, Log)

	if err != nil {
		t.Fatalf("NewProducer fail. | err: %s", err)
	}

	go produceMessage(producer)

	consumerGroupConf := &ConsumerGroupConfig{
		Brokers:         []string{"localhost:9092"},
		ConsumerConfigs: make([]*ConsumerConfig, 0, 8),
	}

	consumerGroupConf.ConsumerConfigs = append(consumerGroupConf.ConsumerConfigs, &ConsumerConfig{
		Brokers: consumerGroupConf.Brokers,
		GroupId: "consumer_group_1",
		Topic:   "testgo",
		Worker:  10,
	})

	doMessage := func(v []byte) error {
		Log.Infof("receive msg. | value: %s", string(v))

		return nil
	}

	go func() {
		consumer, _ := NewConsumer(consumerGroupConf.ConsumerConfigs[0], Log, doMessage)

		time.Sleep(15 * time.Second)
		consumer.Stop()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	Log.Infof("test over")
}

func produceMessage(p *Producer) {
	for i := 0; i < 100; i++ {
		go p.Send("testgo", strconv.Itoa(i), "")
		time.Sleep(100 * time.Millisecond)
	}
	return
}
