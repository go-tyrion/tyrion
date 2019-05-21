package kafka

import (
	"../../config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	kw *kafka.Writer
}

func NewKakfaProducer(conf string) *Producer {
	p := new(Producer)
	p.kw = kafka.NewWriter(p.resolveConf(conf))

	return p
}

func (p *Producer) resolveConf(conf string) kafka.WriterConfig {
	wc := kafka.WriterConfig{
		Brokers: config.Strings("broker", conf, ","),
		Topic:   config.String("topic", conf),
	}

	return wc
}
