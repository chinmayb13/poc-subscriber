package reader

import (
	"context"
	"poc-subscriber/internal/dao/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPullMessages(t *testing.T) {
	tests := []struct {
		name string
		dbInsertErr error
	}{
		{
			name: "Bad Request Error", 
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			logger, _ := zap.NewProduction()
			mockDB := mocks.DBService{}
			mockDB.On("Add", "message").Return(tt.dbInsertErr)
			clientConofig := SubscriberConfig{
				ProjectID: "project",
				SubID:     "sub",
				Logger:    logger,
				Dao:       &mockDB,
			}

			client := GetClientAndSubscription(ctx, &clientConofig)

			err := client.PullMessages(ctx)
			assert.NotNil(t,err)

		})
	}
}
