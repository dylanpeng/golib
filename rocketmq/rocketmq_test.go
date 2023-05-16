package rocketmq

import (
	"fmt"
	"github.com/dylanpeng/golib/logger"
	"testing"
	"time"
)

func TestRocketmqProducerConsumer(t *testing.T) {
	Log, err := logger.NewLogger(&logger.Config{
		FilePath:   "./logs/confuse",
		Level:      "debug",
		TimeFormat: "2006-01-02 15:04:05.000",
		MaxAgeDay:  30,
	})

	if err != nil {
		t.Fatalf("NewLogger fail. | err: %s", err)
	}

	consumerConf := &ConsumerConfig{
		ConnectConfig: &ConnectConfig{
			Endpoints:     []string{"127.0.0.1:9876"},
			AccessKey:     "",
			SecretKey:     "",
			SecurityToken: "",
		},
		Topic:     "test",
		Group:     "consumer-1",
		Orderly:   false,
		FromFirst: false,
		Tags:      nil,
		Worker:    1,
	}

	doMessage := func(v []byte) error {
		Log.Infof("receive msg. | value: %s", string(v))
		return nil
	}

	go func() {
		consumer, _ := NewConsumer(consumerConf, doMessage, Log)

		time.Sleep(15 * time.Second)
		consumer.Stop()
	}()

	producerConf := &ProducerConfig{
		ConnectConfig: &ConnectConfig{
			Endpoints:     []string{"127.0.0.1:9876"},
			AccessKey:     "",
			SecretKey:     "",
			SecurityToken: "",
		},
		Group:      "",
		RetryTimes: 2,
	}

	producer, err := NewProducer(producerConf, Log)

	if err != nil {
		t.Fatalf("NewProducer fail. | err: %s", err)
	}

	go produceMessage(producer)
	//produceDelayMessage(producer)

	//wg := &sync.WaitGroup{}
	//wg.Add(1)
	//wg.Wait()
	time.Sleep(20 * time.Second)
	Log.Infof("test over")
}

func produceMessage(p *Producer) {
	for i := 0; i < 10; i++ {
		msg := &Message{
			Topic:       "test",
			Payload:     []byte(fmt.Sprintf("%s｜%d", time.Now(), i)),
			Tag:         "",
			Keys:        nil,
			ShardingKey: "",
			DelayLevel:  0,
		}
		p.SendSync(msg)
		//time.Sleep(100 * time.Millisecond)
	}
	return
}

func produceDelayMessage(p *Producer) {

	msg := &Message{
		Topic:       "test",
		Payload:     []byte("delay 10s"),
		Tag:         "",
		Keys:        nil,
		ShardingKey: "",
		// 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h(delay level starts from 1.)
		DelayLevel: 3,
	}
	p.SendSync(msg)

	msg = &Message{
		Topic:       "test",
		Payload:     []byte("delay 5s"),
		Tag:         "",
		Keys:        nil,
		ShardingKey: "",
		// 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h(delay level starts from 1.)
		DelayLevel: 2,
	}
	p.SendSync(msg)
	//time.Sleep(100 * time.Millisecond)

	return
}
