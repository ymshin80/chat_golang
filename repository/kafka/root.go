package kafka

import (
	"chat_server_golang/config"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)
type Kafka struct {
	cfg *config.Config
	producer *kafka.Producer
}


func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error
	log.Println("URL=",cfg.Kafka.URL )
	log.Println("ClientID=",cfg.Kafka.ClientID )

	if k.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URL,
		"client.id": cfg.Kafka.ClientID,
		"acks": "all",
	}); err != nil  {
		return nil , err
	} else {
		return k , nil
	}
}
