package uow

import (
	sqliteRepository "github.com/Mateus-Lacerda/better-mem/internal/database/sqlite/repository"
	"github.com/Mateus-Lacerda/better-mem/internal/repository"
	"github.com/Mateus-Lacerda/better-mem/internal/uow"
	"context"

	"gorm.io/gorm"
)

type SQLiteUnitOfWork[T any, C any] struct {
	*gorm.DB
}

// Do implements uow.UnitOfWork.
func (s *SQLiteUnitOfWork[T, C]) Do(
	ctx context.Context, fn func(repos repository.AllRepositories) (T, error),
) (T, error) {
	var result T

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(s.Repositories(tx))
		return err
	})

	if err != nil {
		return *new(T), err
	}
	return result, nil
}

// Repositories implements uow.UnitOfWork.
func (m *SQLiteUnitOfWork[T, C]) Repositories(tx any) repository.AllRepositories {
	gormTx := tx.(*gorm.DB)
	return repository.AllRepositories{
		Chat:            sqliteRepository.ChatRepositoryWithTransaction(gormTx),
		ShortTermMemory: sqliteRepository.ShortTermMemoryRepositoryWithTransaction(gormTx),
		LongTermMemory:  sqliteRepository.LongTermMemoryRepositoryWithTransaction(gormTx),
	}
}

func NewUnitOfWork[T any, C any](
	client *gorm.DB,
) *SQLiteUnitOfWork[T, any] {
	return &SQLiteUnitOfWork[T, any]{
		client,
	}
}

var _ uow.UnitOfWork[any, any] = (*SQLiteUnitOfWork[any, *gorm.DB])(nil)
