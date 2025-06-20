package service

import (
	"context"
	"fmt"
	"practical-assignment/internal/model"
	"time"
)

// processTask обрабатывает задачу асинхронно, имитируя I/O bound операцию.
// Возвращает ошибку, если обработка не удалась.
func (s *TaskService) processTask(ctx context.Context, taskID string) error {
	if err := validateUUID(taskID); err != nil {
		s.logger.Error("Invalid task id", "task_id", taskID, "error", err)
		return err
	}

	task, err := s.storage.Get(taskID)
	if err != nil {
		s.logger.Warn("Task not found", "task_id", taskID)
		return fmt.Errorf("task not found: %w", err)
	}

	if task.Status != model.StatusPending {
		s.logger.Warn("Invalid task status for processing", "task_id", taskID, "status", string(task.Status))
		return fmt.Errorf("invalid task status: %s", string(task.Status))
	}

	task.Status = model.StatusRunning
	if err := s.validateAndSave(task, "running"); err != nil {
		return err
	}
	s.logger.Info("Task started", "task_id", task.ID)

	extraDuration := time.Duration(s.rand.Intn(3)) * time.Minute

	select {
	case <-time.After(s.taskDuration + extraDuration):
		return s.finalizeTask(task, model.StatusCompleted, "Task completed")
	case <-ctx.Done():
		return s.finalizeTask(task, model.StatusFailed, fmt.Sprintf("Task cancelled: %s", ctx.Err()))
	}
}
