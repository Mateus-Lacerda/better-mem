package uow

import (
	"better-mem/internal/repository"
	"context"
)

type UnitOfWork[T any, C any] interface {
	Repositories(tx C) repository.AllRepositories

	Do(
		ctx context.Context, fn func(
			repos repository.AllRepositories,
		) (T, error),
	) (T, error)
}
