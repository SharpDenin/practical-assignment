package storage

import (
	"practical-assignment/internal/model"
	"sync"
)

//TODO Добавить обработку ошибок, переработать Save/Get/Delete (?)

// InMemoryStorage - структура хранилища
type InMemoryStorage struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

// NewInMemoryStorage - создание нового хранилища
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[string]*model.Task),
	}
}

// Save - сохранение задачи в хранилище
func (s *InMemoryStorage) Save(task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryStorage) Get(id string) (*model.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *InMemoryStorage) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[id]; exists {
		delete(s.tasks, id)
		return true
	}
	return false
}
