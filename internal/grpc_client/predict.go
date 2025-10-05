package protos

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"better-mem/internal/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type PredictClient struct {
	client *grpc.ClientConn
}

var (
	addr = "inference:50051"
	lock = &sync.Mutex{}
	client *PredictClient
)

const (
	predictionTimeout = 180
)


func GetPredictClient() *PredictClient {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		if client == nil {
			grpcClient, err := grpc.NewClient(
				addr,
				grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				),
			)

			if err != nil {
				slog.Error("failed to create grpc client", "err", err)
				panic(err)
			}

			client = &PredictClient{
				client: grpcClient,
			}
		}
	}
	return client
}


func Predict(message string, chatId string, withEmbedding bool) (*core.LabeledMessage, error) {

	predictionClient := NewPredictionClient(GetPredictClient().client)
	ctx, cancel := context.WithTimeout(context.Background(), predictionTimeout*time.Second)

	defer cancel()

	response, err := predictionClient.Predict(ctx, &PredictionRequest{
		Message: message, ReturnEmbedding: withEmbedding,
	})
	if err != nil {
		return nil, err
	}

	labeledMessage := &core.LabeledMessage{
		MessageEmbedding: response.Embedding,
		Label: core.MemoryTypeEnum(response.Label),
		NewMessage: core.NewMessage{
			Message: message,
			ChatId: chatId,
		},
	}

	return labeledMessage, nil
}

func Embed(message string) ([]float32, error) {
	predictionClient := NewPredictionClient(GetPredictClient().client)
	ctx, cancel := context.WithTimeout(context.Background(), predictionTimeout*time.Second)

	defer cancel()

	response, err := predictionClient.Embed(ctx, &EmbedRequest{Message: message})
	if err != nil {
		return nil, err
	}

	return response.Embedding, nil
}
