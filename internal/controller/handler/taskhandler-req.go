package handler

import "cybertask/model"

// Generalized error on Task operations.
type TaskError struct {
	Msg string `json:"msg"`
}

type (
	// ->> taskhandler.Get .
	GetTaskRequest struct {
	}
	// taskhandler.Get ->> .
	GetTaskResponse struct {
		Task model.Task `json:"task"`
	}
)
type (
	// ->> taskhandler.Update .
	UpdateTaskRequest struct {
		Task model.Task `json:"task"`
	}
	// taskhandler.Update ->> .
	UpdateTaskResponse struct {
		Msg string `json:"msg"`
	}
)
type (
	// ->> taskhandler.Delete .
	DeleteTaskRequest struct {
	}
	// taskhandler.Delete ->> .
	DeleteTaskResponse struct {
		Msg string `json:"msg"`
	}
)
type (
	// ->> taskhandler.Create .
	CreateTaskRequest struct {
		Task model.Task `json:"task"`
	}
	// taskhandler.Create ->> .
	CreateTaskResponse struct {
		Msg string `json:"msg"`
	}
)
