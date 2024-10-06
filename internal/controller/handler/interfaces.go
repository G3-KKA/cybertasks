package handler

import (
	"context"
	"cybertask/model"
)

//go:generate mockery --filename=mock_taskusecase.go --name=TaskUsecase --dir=. --structname=MockTaskUsecase --outpkg=mock_handler
type TaskUsecase interface {
	CreateTask(ctx context.Context, task model.Task) error
	UpdateTask(ctx context.Context, task model.Task) error
	GetTask(ctx context.Context, id model.TaskID) (model.Task, error)
	DeleteTask(ctx context.Context, id model.TaskID) error
}
