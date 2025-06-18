package model

import "time"

// TaskStatus - возможные статусы задачи
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

// Task - структура I/O bound задачи
type Task struct {
	ID        string     `json:"id"`
	Status    TaskStatus `json:"task_status"`
	CreatedAt time.Time  `json:"created_at"`
	//UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Result      string     `json:"result"`
}
