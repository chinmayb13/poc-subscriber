package reader

import (
	"context"
	"log"
	"poc-subscriber/internal/dao"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type SubscriberService interface {
	PullMessages(ctx context.Context) error
	CloseClient()
}

type SubscriberConfig struct {
	ProjectID string
	SubID     string
	Logger    *zap.Logger
	Dao       dao.DBService
}

type psClient struct {
	*pubsub.Client
	logger       *zap.Logger
	subscription *pubsub.Subscription
	dao          dao.DBService
}

func GetClientAndSubscription(ctx context.Context, config *SubscriberConfig) SubscriberService {
	client, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("client creation failed %s", err.Error())
	}
	sub := client.Subscription(config.SubID)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10
	sub.ReceiveSettings.MaxOutstandingBytes = 1e10

	return &psClient{
		Client:       client,
		logger:       config.Logger,
		subscription: sub,
		dao:          config.Dao,
	}
}

func (client *psClient) CloseClient() {
	client.Close()
}

func (client *psClient) PullMessages(ctx context.Context) error {
	client.logger.Info("listening to pubsub...")
	// msgChan := make(chan *pubsub.Message)
	// errChan := make(chan error)
	// go func() {
	// 	for msg := range msgChan {

	// 	}
	// }()
	// defer close(msgChan)
	cctx, cancel := context.WithCancel(ctx)
	var err, receiveErr error

	receiveErr = client.subscription.Receive(cctx, func(c context.Context, m *pubsub.Message) {
		client.logger.Info("received pubsub message", zap.String("msg", string(m.Data)))
		err = client.dao.AddRow(string(m.Data))
		if err != nil {
			client.logger.Error("db insert failed", zap.String("err", err.Error()))
			cancel()
		}
		m.Ack()
	})
	if err != nil {
		receiveErr = err
	}

	return receiveErr
}
