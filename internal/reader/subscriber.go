package reader

import (
	"context"
	"log"
	"poc-subscriber/internal/dao"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaService interface {
	ReadFromKafka()
}

func GetKafkaService(dao dao.DBService, kafkaconfig KafkaConfig) KafkaService {
	return &serviceClient{
		kafkaReader: kafkaconfig.getReader(),
		dao:         dao,
	}
}

type KafkaConfig struct {
	Topic     string
	Broker    string
	GroupID   string
	Partition int
}

type serviceClient struct {
	kafkaReader *kafka.Reader
	dao         dao.DBService
}

func (c *KafkaConfig) getReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.Broker},
		Topic:     c.Topic,
		Partition: c.Partition,
		GroupID:   c.GroupID,
	})
}

func (client *serviceClient) ReadFromKafka() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()
	ctx := context.Background()
	for {
		msg, err := client.kafkaReader.ReadMessage(ctx)
		if err != nil {
			logger.Info("kafka read failed", zap.String("err", err.Error()))
			break
		}
		logger.Info("received kafka message", zap.String("val", string(msg.Value)))

		err = client.dao.AddRow(string(msg.Value))
		if err != nil {
			logger.Error("db insert failed", zap.String("err", err.Error()))
			break
		}
		err = client.kafkaReader.CommitMessages(ctx, msg)
		if err != nil {
			logger.Error("message commit failed", zap.String("err", err.Error()))
			break
		}
	}

}
