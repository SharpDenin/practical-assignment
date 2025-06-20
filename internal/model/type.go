package model

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// MaxResultLength - Максимальная длина результата задачи
const MaxResultLength = 1024

// TaskStatus - тип, представляющий возможные статусы задачи.
// Значения определяются константами StatusPending, StatusRunning, StatusCompleted, StatusFailed.
type TaskStatus string

// Константы статусов задачи
const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// Tasker - интерфейс для задач, поддерживающий разные типы задач в будущем.
type Tasker interface {
	GetID() string
	GetStatus() TaskStatus
	Validate() error
}

// Task - структура, представляющая I/O bound задачу.
// Используется для хранения и передачи данных о задаче через API и хранилище.
type Task struct {
	ID          string     `json:"id"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Result      string     `json:"result,omitempty"`
}

// ID возвращает идентификатор задачи, реализуя интерфейс Tasker
func (t *Task) GetID() string {
	return t.ID
}

// Status возвращает статус задачи, реализуя интерфейс Tasker
func (t *Task) GetStatus() TaskStatus {
	return t.Status
}

// Validate проверяет валидность полей задачи.
// Возвращает ошибку, если:
// - ID не является валидным UUID
// - Status невалидный
// - CreatedAt является нулевым (time.Time{})
// - CompletedAt не nil для незавершённых задач (pending, running)
// - Result превышает максимальную длину
func (t *Task) Validate() error {
	if _, err := uuid.Parse(t.ID); err != nil {
		return errors.New("invalid task id: must be a valid UUID")
	}

	if !t.Status.IsValid() {
		return errors.New("invalid task status")
	}

	if t.CreatedAt.IsZero() {
		return errors.New("created_at must not be zero")
	}

	if t.CompletedAt != nil && t.Status != StatusCompleted && t.Status != StatusFailed {
		return errors.New("completed_at must be nil for pending or running tasks")
	}
	if len(t.Result) > MaxResultLength {
		return errors.New("result too long")
	}
	return nil
}
