package main

import (
	"context"
	"extservice/internal/extservice"
	"extservice/model"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	rdata "github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

type service struct {
	tasks map[model.TaskID]*model.Task

	// Every 5 requests will produce an error.
	errcounter atomic.Uint64
	mx         sync.Mutex
	extservice.UnimplementedTaskServiceServer
}

func (srvc *service) newtask() {
	srvc.mx.Lock()
	hdr := rdata.FirstName(-1) + rdata.Address() + "visit"
	id := uuid.New()
	tsk := model.Task{
		Id:          id,
		Header:      hdr,
		Description: "visit" + strconv.Itoa(int(srvc.errcounter.Load())),
		CreatedAt:   time.Now(),
		Status:      false,
	}
	srvc.tasks[id] = &tsk
	defer srvc.mx.Unlock()
}

// GetTasks returns all tasks of service, or error every 5 requests.
func (srvc *service) GetTasks(ctx context.Context, _ *extservice.GetTasksRequest) (*extservice.GetTasksResponse, error) {
	srvc.mx.Lock()
	defer srvc.mx.Unlock()

	// Random internal errors.
	if (srvc.errcounter.Load() % 5) == 0 {
		rsp := &extservice.GetTasksResponse{
			// Oneof.
			Rsp: &extservice.GetTasksResponse_Err{
				Err: &extservice.Error{
					Msg: ErrInernal.Error(),
				},
			},
		}

		return rsp, ErrInernal
	}
	srvc.errcounter.Add(1)

	// Return existing tasks.
	tasks := make([]*extservice.Task, 0, len(srvc.tasks))
	for _, t := range srvc.tasks {
		tasks = append(tasks, mapTaskToGrpc(*t))
	}
	rsp := &extservice.GetTasksResponse{
		// Oneof.
		Rsp: &extservice.GetTasksResponse_Tasks{
			Tasks: &extservice.Tasks{
				Data: tasks,
			},
		},
	}

	return rsp, nil
}

// UpdateTaskStatus update Status field of one specific task,
// or error if not exist,
// or error every 5 requests.
func (srvc *service) UpdateTaskStatus(ctx context.Context, req *extservice.UpdateTaskStatusRequest) (*extservice.UpdateTaskStatusResponse, error) {
	srvc.mx.Lock()
	defer srvc.mx.Unlock()

	// Random internal errors.
	if (srvc.errcounter.Load() % 5) == 0 {
		rsp := &extservice.UpdateTaskStatusResponse{
			// Oneof.
			Rsp: &extservice.UpdateTaskStatusResponse_Err{
				Err: &extservice.Error{
					Msg: ErrInernal.Error(),
				},
			},
		}

		return rsp, ErrInernal
	}
	srvc.errcounter.Add(1)

	uid := uuid.MustParse(req.GetId())
	if _, ok := srvc.tasks[uid]; !ok {
		rsp := &extservice.UpdateTaskStatusResponse{
			Rsp: &extservice.UpdateTaskStatusResponse_Err{
				Err: &extservice.Error{
					Msg: ErrTaskNotExist.Error(),
				},
			},
		}

		return rsp, ErrTaskNotExist
	}

	srvc.tasks[uid].Status = req.GetStatus()

	rsp := &extservice.UpdateTaskStatusResponse{
		Rsp: &extservice.UpdateTaskStatusResponse_Msg{
			Msg: "updated successfuly",
		},
	}

	return rsp, nil
}
