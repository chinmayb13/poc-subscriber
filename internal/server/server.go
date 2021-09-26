package server

import (
	"log"
	"poc-subscriber/internal/dao"
	"poc-subscriber/internal/reader"
	"sync"

	"go.uber.org/zap"
)

func RunServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()

	ysqlConfig := dao.YsqlConfig{
		Host:     "172.31.99.57",
		Port:     5433,
		User:     "yugabyte",
		Password: "yugabyte",
		DbName:   "poc_demo",
	}
	dao := dao.GetDBService(&ysqlConfig)
	config := reader.KafkaConfig{
		Topic:   "poc-topic",
		Broker:  "172.31.99.57:9092",
		GroupID: "poc-group",
	}

	kafkaService := reader.GetKafkaService(dao, config)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		kafkaService.ReadFromKafka()
		wg.Done()
	}()
	logger.Info("starting service ...")
	wg.Wait()

}
