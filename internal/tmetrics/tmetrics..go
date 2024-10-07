package tmetrics

import (
	"context"
	"cybertask/config"
	"cybertask/internal/logger"
	"cybertask/model"
	"cybertask/pkg/streampool"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type tmetricsService struct {
	pool     *streampool.Pool
	tasks    map[model.TaskID]*externalService
	external []*externalService

	shutdown sync.Once
	mx       sync.RWMutex
}

// GetTasks implements ucase.TaskMetrics.
func (srvc *tmetricsService) GetTasks(ctx context.Context) ([]model.Task, error) {

	srvc.mx.Lock()
	defer srvc.mx.Unlock()

	// Cheap temporal storage for slice headers.
	perExtServiceTasks := make([][]model.Task, len(srvc.external))

	for idx, ext := range srvc.external {

		perExtServiceTasks[idx] = ext.flushCache()

		// Caching all the tasks.
		// TODO: map will only grow, never cleaned.
		for _, t := range perExtServiceTasks[idx] {
			srvc.tasks[t.ID] = ext
		}

	}

	tasks := slices.Concat(perExtServiceTasks...)

	return tasks, nil
}

// UpdateTaskStatus implements ucase.TaskMetrics.
func (srvc *tmetricsService) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status bool) error {

	srvc.mx.Lock()
	defer srvc.mx.Unlock()

	ext, exist := srvc.tasks[id]
	if !exist {
		return ErrTaskDoesNotExist
	}
	return ext.updateTaskStatus(ctx, id, status)

}

// TaskMetrics constructor.
func New(l *logger.Logger, cfg config.TaskMetrics) (*tmetricsService, error) {
	service := &tmetricsService{
		pool:     streampool.NewStreamPool(),
		tasks:    make(map[uuid.UUID]*externalService, cfg.PerServiceTableSize*uint(len(cfg.ExtServices))),
		external: make([]*externalService, 0, len(cfg.ExtServices)*2),
		shutdown: sync.Once{},
		mx:       sync.RWMutex{},
	}

	service.mx.Lock()
	defer service.mx.Unlock()

	start := make(chan struct{})
	defer close(start)

	for _, eservice := range cfg.ExtServices {
		rawClient, err := grpc.NewClient(eservice.Address)

		// Retry connection attempts.
		if err != nil {
			for range cfg.ConnectionReties {
				time.Sleep(cfg.RetryAfter)
				if rawClient, err = grpc.NewClient(eservice.Address); err == nil {
					break
				}
			}
			if err != nil {
				return nil, err
			}
		}

		taskClient := NewTaskServiceClient(rawClient)
		es := externalService{
			client: taskClient,
			l:      l,
			cached: make([]model.Task, 0),
			errs:   make([]error, 0),
			mx:     sync.RWMutex{},
		}
		err = service.pool.Go(eservice.Name, es.toPool(start, eservice.Autoupdate, rawClient.Close))

		service.external = append(service.external, &es)

		if err != nil {
			return nil, err
		}

	}

	return service, nil
}
func (srvc *tmetricsService) Shutdown(ctx context.Context) error {
	srvc.shutdown.Do(func() {
		srvc.pool.ShutdownWait()
	})

	return nil
}
