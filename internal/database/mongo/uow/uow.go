package uow

import (
	"context"
	"better-mem/internal/repository"
	"better-mem/internal/uow"
	mongoRepository "better-mem/internal/database/mongo/repository"
	"better-mem/internal/database/mongo"

	mongoClient "go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork[T any] struct {
	*mongoClient.Database
}

// Do implements uow.UnitOfWork.
func (m *MongoUnitOfWork[T]) Do(
	ctx context.Context, fn func(repos repository.AllRepositories)(T, error),
) (T, error) {
	transactionResult := *new(T)
	sess, err := m.Database.Client().StartSession()
	if err != nil {
		return transactionResult, err
	}
	defer func() {
		sess.AbortTransaction(ctx)
		sess.EndSession(ctx)
	}()

	result, err := sess.WithTransaction(
		ctx,
		func(sessCtx mongoClient.SessionContext) (any, error) {
			return fn(m.Repositories())
		},
	)
	if err != nil {
		return transactionResult, err
	}
	return result.(T), nil
}

// Repositories implements uow.UnitOfWork.
func (m *MongoUnitOfWork[T]) Repositories() repository.AllRepositories {
	return repository.AllRepositories{
		Chat: mongoRepository.NewChatRepository(),
		ShortTermMemory: mongoRepository.NewShortTermMemoryRepository(),
		LongTermMemory: mongoRepository.NewLongTermMemoryRepository(),
	}
}

func NewUnitOfWork[T any](
	client *mongo.MongoClient,
) *MongoUnitOfWork[T] {
	return &MongoUnitOfWork[T]{
		Database: mongo.GetMongoDatabase(),
	}
}

var _ uow.UnitOfWork[any] = (*MongoUnitOfWork[any])(nil)
