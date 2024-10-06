package main

import (
	"extservice/internal/extservice"
	"extservice/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapTaskToGrpc(t model.Task) *extservice.Task {
	return &extservice.Task{
		Id:          t.Id,
		Header:      t.Header,
		Description: t.Description,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		Status:      t.Status,
	}
}

func mapGrpcToTask(t *extservice.Task) model.Task {
	return model.Task{
		Id:          t.Id,
		Header:      t.Header,
		Description: t.Description,
		CreatedAt:   t.CreatedAt.AsTime(),
		Status:      t.Status,
	}
}
