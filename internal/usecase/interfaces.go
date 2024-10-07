package ucase

import (
	"context"
	"cybertask/model"
)

//go:generate mockery --filename=mock_task_usecase.go --name=TaskUsecase --dir=. --structname=MockTaskUsecase --outpkg=mock_usecase
type TaskRepo interface {
	CreateTask(ctx context.Context, task model.Task) error
	UpdateTask(ctx context.Context, task model.Task) error
	GetTask(ctx context.Context, id model.TaskID) (model.Task, error)
	DeleteTask(ctx context.Context, id model.TaskID) error
	CreateTasks(ctx context.Context, tasks []model.Task) error
}

//go:generate mockery --filename=mock_task_metrics.go --name=TaskMetrics --dir=. --structname=MockTaskMetrics --outpkg=mock_usecase
type TaskMetrics interface {
	GetTasks(ctx context.Context) ([]model.Task, error)
	UpdateTaskStatus(ctx context.Context, id model.TaskID, status model.TaskStatus) error
}
