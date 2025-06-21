package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"practical-assignment/internal/model"
	"practical-assignment/internal/storage"
)

// TaskProcessor - интерфейс для обработки задач.
// Поддерживает создание, получение, удаление и список задач.
type TaskProcessor interface {
	CreateTask(ctx context.Context) (string, error)
	GetTask(id string) (*model.Task, error)
	DeleteTask(id string) error
	ListTasks() ([]*model.Task, error)
}

// TaskService - сервис для управления задачами.
// Реализует бизнес-логику создания, обработки и управления задачами.
type TaskService struct {
	storage      storage.Storage
	logger       *slog.Logger
	taskDuration time.Duration
	rand         *rand.Rand
}

// NewTaskService создаёт новый сервис задач.
// Принимает хранилище и логгер.
func NewTaskService(storage storage.Storage, logger *slog.Logger) *TaskService {
	return &TaskService{
		storage:      storage,
		logger:       logger,
		taskDuration: time.Minute * 3,
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateTask создаёт новую задачу и запускает её обработку в горутине.
// Возвращает ID задачи или ошибку при сохранении.
func (s *TaskService) CreateTask(ctx context.Context) (string, error) {
	task := &model.Task{
		ID:        uuid.New().String(),
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.validateAndSave(task, "create"); err != nil {
		return "", err
	}

	go func() {
		_ = s.processTask(context.Background(), task.ID)
	}()

	return task.ID, nil
}

// GetTask возвращает задачу по ID.
// Возвращает ошибку, если задача не найдена или ID невалиден.
func (s *TaskService) GetTask(id string) (*model.Task, error) {
	if err := validateUUID(id); err != nil {
		s.logger.Error("Invalid task id", "task_id", id, "error", err)
		return nil, err
	}
	task, err := s.storage.Get(id)
	if err != nil {
		s.logger.Warn("Task not found", "task_id", id, "error", err)
		return nil, fmt.Errorf("task not found: %w", err)
	}

	s.logger.Info("Task retrieved", "task_id", task.ID)
	return task, nil
}

// DeleteTask удаляет задачу по ID.
// Возвращает ошибку, если задача не найдена или ID невалиден.
func (s *TaskService) DeleteTask(id string) error {
	if err := validateUUID(id); err != nil {
		s.logger.Error("Invalid task id", "task_id", id, "error", err)
		return err
	}
	if err := s.storage.Delete(id); err != nil {
		s.logger.Error("Failed to delete task", "task_id", id, "error", err)
		return fmt.Errorf("failed to delete task: %w", err)
	}
	s.logger.Info("Task deleted", "task_id", id)
	return nil
}

// ListTasks возвращает список всех задач.
func (s *TaskService) ListTasks() ([]*model.Task, error) {
	tasks, err := s.storage.List()
	if err != nil {
		s.logger.Error("Failed to list tasks", "error", err)
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	s.logger.Info("Tasks listed", "count", len(tasks))
	return tasks, nil
}

func (s *TaskService) validateAndSave(task *model.Task, stage string) error {
	if err := task.Validate(); err != nil {
		s.logger.Error("Invalid task", "task_id", task.ID, "stage", stage, "error", err)
		return fmt.Errorf("invalid task: %w", err)
	}
	if err := s.storage.Save(task); err != nil {
		s.logger.Error("Failed to save task", "task_id", task.ID, "stage", stage, "error", err)
		return fmt.Errorf("failed to save task: %w", err)
	}
	return nil
}

func (s *TaskService) finalizeTask(task *model.Task, status model.TaskStatus, result string) error {
	task.Status = status
	task.Result = result
	if status == model.StatusCompleted || status == model.StatusFailed {
		now := time.Now()
		task.CompletedAt = &now
	}
	return s.validateAndSave(task, string(status))
}

func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return nil
}
