package storage

import (
	"errors"
	"github.com/google/uuid"
	"practical-assignment/internal/model"
	"sync"
)

// Storage - интерфейс для хранилища задач.
// Поддерживает операции сохранения, получения, удаления и списка задач.
type Storage interface {
	Save(task *model.Task) error
	Get(id string) (*model.Task, error)
	Delete(id string) error
	List() ([]*model.Task, error)
}

// InMemoryStorage - in-memory хранилище для задач.
// Использует map с конкурентным доступом через RWMutex.
type InMemoryStorage struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

// NewInMemoryStorage создаёт новое in-memory хранилище для задач.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[string]*model.Task),
	}
}

// Save сохраняет задачу в хранилище.
// Проверяет валидность задачи и ID.
func (s *InMemoryStorage) Save(task *model.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	if _, err := uuid.Parse(task.ID); err != nil {
		return errors.New("invalid task id: must be a valid UUID")
	}
	if err := task.Validate(); err != nil {
		return errors.New("invalid task status" + err.Error())
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

// Get возвращает задачу по ID.
// Возвращает ошибку, если ID невалиден или задача не найдена.
func (s *InMemoryStorage) Get(id string) (*model.Task, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid task id: must be a valid UUID")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// Delete удаляет задачу по ID.
// Возвращает ошибку, если ID невалиден или задача не найдена.
func (s *InMemoryStorage) Delete(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid task id: must be a valid UUID")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}

// List возвращает список всех задач в хранилище.
func (s *InMemoryStorage) List() ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
