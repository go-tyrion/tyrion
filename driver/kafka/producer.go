package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"lib/config"
	e "lib/error"
	"time"
)

func NewKakfaProducer(conf string) *Producer {
	p := &Producer{
		conf:    conf,
		writers: make(map[string]*kafka.Writer),
	}

	p.init()

	return p
}

type Producer struct {
	conf string

	writers map[string]*kafka.Writer
}

func (p *Producer) init() {
	topics := config.Strings("topics", p.conf, ",")
	for _, topic := range topics {
		wc := p.resolveConf(topic, config.Strings("brokers", p.conf, ","))
		p.writers[topic] = kafka.NewWriter(wc)
	}
}

// 生产消息
func (p *Producer) Product(topic string, msg []byte) error {
	if w, ok := p.writers[topic]; ok {
		return w.WriteMessages(context.Background(), kafka.Message{
			Topic: topic,
			Key:   p.defaultKey(),
			Value: msg,
		})
	}

	return e.New("no such topic")
}

func (p *Producer) defaultKey() []byte {
	return []byte("tyrion-micro-fw")
}

func (p *Producer) resolveConf(topic string, brokers []string) kafka.WriterConfig {
	wc := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.RoundRobin{},
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		RequiredAcks: 1,
		Async:        false, // 使用同步的方式，虽然相对较慢，但是安全
	}

	return wc
}
