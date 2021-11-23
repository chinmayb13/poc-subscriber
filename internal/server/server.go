package server

import (
	"context"
	"log"
	"poc-subscriber/config"
	"poc-subscriber/internal/dao"
	"poc-subscriber/internal/reader"

	"go.uber.org/zap"
)

func RunServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
	}
	defer logger.Sync()

	appConfig, err := config.LoadConfig("../..")
	if err != nil {
		logger.Fatal("error reading env file", zap.String("err", err.Error()))
	}

	ctx := context.Background()

	ysqlConfig := dao.YsqlConfig{
		Host:     appConfig.DB.Host,
		Port:     appConfig.DB.Port,
		User:     appConfig.DB.User,
		Password: appConfig.DB.Password,
		DbName:   appConfig.DB.Name,
	}
	dao := dao.GetDBService(&ysqlConfig)
	// config := reader.KafkaConfig{
	// 	Topic:   appConfig.Kafka.Topic,
	// 	Broker:  appConfig.Kafka.Broker,
	// 	GroupID: appConfig.Kafka.GroupID,
	// }

	// kafkaService := reader.GetKafkaService(dao, config)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	kafkaService.ReadFromKafka()
	// 	wg.Done()
	// }()
	// logger.Info("starting service ...")
	// wg.Wait()
	pubsubCfg := &reader.SubscriberConfig{
		ProjectID: appConfig.PubSub.ProjectID,
		SubID:     appConfig.PubSub.SubID,
		Logger:    logger,
		Dao:       dao,
	}
	pubsubService := reader.GetClientAndSubscription(ctx, pubsubCfg)
	defer pubsubService.CloseClient()
	logger.Info("starting poc-subscriber...")
	err = pubsubService.PullMessages(ctx)
	if err != nil {
		logger.Fatal("pubsub subscriber failed", zap.String("err", err.Error()))
	}

}
