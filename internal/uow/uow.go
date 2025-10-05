package uow

import (
	"context"
	"better-mem/internal/repository"
)

type UnitOfWork[T any] interface {
	Repositories() repository.AllRepositories

	Do(
		ctx context.Context, fn func(
			repos repository.AllRepositories,
		) (T, error),
	) (T, error)
}
