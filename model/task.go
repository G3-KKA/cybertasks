package model

import (
	"time"
)

type (
	TaskID      = uint64
	TaskHeader  = string
	TaskDesc    = string
	TaskCreated = time.Time
	TaskStatus  = bool
)

type Task struct {
	ID          TaskID      `json:"id"`
	Header      TaskHeader  `json:"header" validate:"min=1,max=255"`
	Description TaskDesc    `json:"description,omitempty"`
	CreatedAt   TaskCreated `json:"created_at" validate:"required"`
	Status      TaskStatus  `json:"status"`
}
