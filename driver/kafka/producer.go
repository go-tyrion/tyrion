package kafka

import (
	"github.com/Shopify/sarama"
	"lib/config"
)

type KafkaProduceer struct {
	conf     string
	producer sarama.AsyncProducer
}

func (p *KafkaProduceer) Init(conf string) error {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(config.Strings("brokers", conf, ","), c)
	if err != nil {
		return err
	}

	p.producer = producer

	return nil
}

func (p *KafkaProduceer) Produce(topic string, content []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: "",
		Key:   sarama.StringEncoder(""),
		Value: sarama.ByteEncoder(""),
	}
	p.producer.Input() <- msg

	return nil
}
