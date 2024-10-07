package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	TaskID      = uuid.UUID
	TaskHeader  = string
	TaskDesc    = string
	TaskCreated = time.Time
	TaskStatus  = bool
)

type Task struct {
	ID          TaskID      `json:"id" `
	Header      TaskHeader  `json:"header" validate:"min=1,max=255"`
	Description TaskDesc    `json:"description,omitempty"  bun:",nullzero"`
	CreatedAt   TaskCreated `json:"created_at" validate:"required" bun:"created_at"`
	Status      TaskStatus  `json:"status"`
}
