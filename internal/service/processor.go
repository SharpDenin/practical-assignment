package service

import (
	"context"
	"github.com/google/uuid"
	"math/rand"
	"practical-assignment/internal/model"
	"practical-assignment/internal/storage"
	"time"
)

// TaskService - бизнес-логика задач
type TaskService struct {
	storage *storage.InMemoryStorage
}

// NewTaskService - создание нового Таск-сервиса
func NewTaskService(storage *storage.InMemoryStorage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) CreateTask(ctx context.Context) (string, error) {
	task := &model.Task{
		ID:        uuid.New().String(),
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.storage.Save(task); err != nil {
		return "", err
	}
	// Запускаем обработку таски в горутине
	go s.processTask(context.Background(), task.ID)

	return task.ID, nil
}

// processTask simulates a long-running I/O bound task
func (s *TaskService) processTask(ctx context.Context, taskID string) {
	task, exists := s.storage.Get(taskID)
	if !exists {
		return
	}

	task.Status = model.StatusRunning
	err := s.storage.Save(task)
	if err != nil {
		return
	}

	// Simulate I/O bound task (3-5 minutes)
	select {
	case <-time.After(time.Duration(3+rand.Intn(3)) * time.Minute):
		task.Status = model.StatusCompleted
		completedAt := time.Now()
		task.CompletedAt = &completedAt
		task.Result = "Task completed"
		err := s.storage.Save(task)
		if err != nil {
			return
		}
	case <-ctx.Done():
		task.Status = model.StatusFailed
		task.Result = "Task cancelled"
		err := s.storage.Save(task)
		if err != nil {
			return
		}
	}
}

func (s *TaskService) GetTask(id string) (*model.Task, bool) {
	task, exists := s.storage.Get(id)
	if !exists {
		return nil, false
	}
	return task, true
}

func (s *TaskService) DeleteTask(id string) bool {
	return s.storage.Delete(id)
}
