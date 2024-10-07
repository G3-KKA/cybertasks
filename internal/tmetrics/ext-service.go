package tmetrics

import (
	context "context"
	"cybertask/internal/logger"
	"cybertask/model"
	"cybertask/pkg/streampool"
	sync "sync"
	"time"

	"github.com/google/uuid"
)

type externalService struct {
	client TaskServiceClient

	l      *logger.Logger
	cached []model.Task
	errs   []error
	mx     sync.RWMutex
}

func (ext *externalService) loadCache() error {

	ext.mx.Lock()
	defer ext.mx.Unlock()

	rsp, err := ext.client.GetTasks(context.TODO(), &GetTasksRequest{})
	if err != nil {
		ext.errs = append(ext.errs, err)

		return err
	}

	tasks := rsp.GetTasks()
	for _, t := range tasks.Data {
		ext.cached = append(ext.cached, mapGrpcToTask(t))
	}

	return nil

}
func (ext *externalService) flushCache() []model.Task {

	ext.mx.Lock()
	defer ext.mx.Unlock()

	tasks := ext.cached
	ext.cached = make([]model.Task, 0, len(ext.cached))

	return tasks

}
func (ext *externalService) toPool(start <-chan struct{}, autoupdate time.Duration, callback func() error) streampool.PoolFunc {
	return streampool.PoolFunc(func(stop <-chan struct{}) {

		ticker := time.NewTicker(autoupdate)
		defer ticker.Stop()
		if callback != nil {
			defer callback()
		}

		<-start
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
			}

			err := ext.loadCache()
			if err != nil {
				ext.l.Err(err).Send()
			}

		}
	})
}
func (ext *externalService) updateTaskStatus(ctx context.Context, id uuid.UUID, status bool) error {
	_, err := ext.client.UpdateTaskStatus(ctx, &UpdateTaskStatusRequest{
		Id:     id.String(),
		Status: status,
	})
	if err != nil {
		return err
	}

	return nil
}
