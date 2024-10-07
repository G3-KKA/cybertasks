package ucase

import (
	"context"
	"cybertask/internal/controller/handler"
	"cybertask/internal/logger"
	"cybertask/model"
	"errors"

	"github.com/google/uuid"
)

var _ handler.TaskUsecase = (*taskUcase)(nil)

type taskUcase struct {
	l *logger.Logger

	repo    TaskRepo
	metrics TaskMetrics
}

func NewTaskUsecase(l *logger.Logger, repo TaskRepo, metrics TaskMetrics) *taskUcase {
	return &taskUcase{
		l:       l,
		repo:    repo,
		metrics: metrics,
	}
}

// GetTask implements handler.TaskUsecase.
func (uc *taskUcase) GetTask(ctx context.Context, id uuid.UUID) (model.Task, error) {

	var (
		errMetrics     error
		errCreateTasks error
	)

	tasks, errMetrics := uc.metrics.GetTasks(ctx)
	if errMetrics == nil {
		errCreateTasks = uc.repo.CreateTasks(ctx, tasks)
	}

	task, err := uc.repo.GetTask(ctx, id)
	if err != nil {
		// Metrics failed, so as database.
		err = errors.Join(err, errMetrics, errCreateTasks)

		return model.Task{}, err
	}

	err = errors.Join(err, errMetrics, errCreateTasks)
	if err != nil {
		// Metrics failed, but task exist in database.
		err = errors.Join(err, ErrNotCritical)

		return task, err
	}

	return task, nil

}

// UpdateTask implements handler.TaskUsecase.
func (uc *taskUcase) UpdateTask(ctx context.Context, task model.Task) error {
	err := uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	err = uc.metrics.UpdateTaskStatus(ctx, task.ID, task.Status)
	if err != nil {
		err = errors.Join(ErrNotCritical, err)

		return err
	}

	return nil
}

// CreateTask implements handler.TaskUsecase.
func (uc *taskUcase) CreateTask(ctx context.Context, task model.Task) error {
	return uc.repo.CreateTask(ctx, task)
}

// DeleteTask implements handler.TaskUsecase.
func (uc *taskUcase) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteTask(ctx, id)
}
