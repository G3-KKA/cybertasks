package tmetrics

import (
	"cybertask/model"

	"github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func mapTaskToGrpc(t model.Task) *Task {
	return &Task{
		Id:          t.ID.String(),
		Header:      t.Header,
		Description: t.Description,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		Status:      t.Status,
	}
}

func mapGrpcToTask(t *Task) model.Task {
	return model.Task{
		ID:          uuid.MustParse(t.Id),
		Header:      t.Header,
		Description: t.Description,
		CreatedAt:   t.CreatedAt.AsTime(),
		Status:      t.Status,
	}
}
