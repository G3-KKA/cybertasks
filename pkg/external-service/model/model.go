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
	Id          TaskID      `json:"id"`
	Header      TaskHeader  `json:"header"`
	Description TaskDesc    `json:"description"`
	CreatedAt   TaskCreated `json:"created_at"`
	Status      TaskStatus  `json:"status"`
}
