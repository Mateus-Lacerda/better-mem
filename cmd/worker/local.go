//go:build local

package main

import (
	"better-mem/internal/config"
	"better-mem/internal/database/sqlite"
	"better-mem/internal/database/sqlite/repository"
	vectorRepo "better-mem/internal/database/sqlite/repository/vector"
	"better-mem/internal/database/sqlite/uow"
	contracts "better-mem/internal/repository"
	vectorContracts "better-mem/internal/repository/vector"
	"better-mem/internal/task"
	"better-mem/internal/task/handler"
	uowContracts "better-mem/internal/uow"
	"context"
	"log"
	"time"

	"github.com/khepin/liteq"
)

type scheduler struct {
	jqueue          *liteq.JobQueue
	intervalSeconds int
	queue           string
	runner          func(context.Context, *liteq.Job) error
}

// Enqueues the job and then runs it
func (s *scheduler) run(ctx context.Context, job *liteq.Job) error {
	s.jqueue.QueueJob(context.Background(), liteq.QueueJobParams{
		Queue:        s.queue,
		ExecuteAfter: time.Now().Add(time.Duration(s.intervalSeconds) * time.Second).Unix(),
		Job:          job.Job,
	})
	return s.runner(ctx, job)
}

func startScheduler() {}

func startConsumer(
	messageHandler *handler.MessageTaskHandler,
	manageShortTermMemoryHandler *handler.MemoryManagementHandler,
) {

	db, err := sqlite.GetDb().DB()
	if err != nil {
		log.Fatal(err)
	}
	jqueue := liteq.New(db)
	scheduler := scheduler{
		jqueue,
		config.MemoryManagement.ManageSTMemoryTaskPeriodInt,
		task.ManageMemoryTaskName,
		manageShortTermMemoryHandler.HandleManageMemory,
	}
	go jqueue.Consume(
		context.Background(),
		liteq.ConsumeParams{
			Queue:             task.ClassifyMessageTaskName,
			VisibilityTimeout: 20,
			Worker:            messageHandler.HandleClassifyMemoryTask,
		},
	)
	go jqueue.Consume(
		context.Background(),
		liteq.ConsumeParams{
			Queue:             task.StoreLongTermMemoryTaskName,
			VisibilityTimeout: 20,
			Worker:            messageHandler.HandleStoreLongTermMemoryTask,
		},
	)
	go jqueue.Consume(
		context.Background(),
		liteq.ConsumeParams{
			Queue:             task.StoreShortTermMemoryTaskName,
			VisibilityTimeout: 20,
			Worker:            messageHandler.HandleStoreShortTermMemoryTask,
		},
	)
	go jqueue.Consume(
		context.Background(),
		liteq.ConsumeParams{
			Queue:             task.ManageMemoryTaskName,
			VisibilityTimeout: 20,
			Worker:            scheduler.run,
		},
	)
	err = jqueue.QueueJob(
		context.Background(),
		liteq.QueueJobParams{
			Queue: task.ManageMemoryTaskName,
			Job:   `{}`,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func getRepositories() (
	contracts.ChatRepository,
	contracts.LongTermMemoryRepository,
	contracts.ShortTermMemoryRepository,
	vectorContracts.MemoryVectorRepository,
	uowContracts.UnitOfWork[int, any],
) {
	chatRepository := repository.NewChatRepository()
	longTermMemoryRepository := repository.NewLongTermMemoryRepository()
	shortTermMemoryRepository := repository.NewShortTermMemoryRepository()
	memoryVectorRepository := vectorRepo.NewMemoryRepository()
	sqliteIntUow := uow.NewUnitOfWork[int, any](sqlite.GetDb())
	sqlite.InitDb()
	sqlite.Migrate(sqlite.GetDb())
	return chatRepository, longTermMemoryRepository, shortTermMemoryRepository, memoryVectorRepository, sqliteIntUow
}
