package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"practical-assignment/internal/model"
	"practical-assignment/internal/service"
)

// TaskHandler - обработчик HTTP-запросов для задач.
// Использует TaskProcessor для бизнес-логики и логирование через slog.
type TaskHandler struct {
	service service.TaskProcessor
	logger  *slog.Logger
}

// NewTaskHandler создаёт новый обработчик задач.
// Принимает сервис и логгер.
func NewTaskHandler(service service.TaskProcessor, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

// errorResponse - структура для JSON-ответов с ошибками.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON пишет JSON-ответ с указанным статусом.
func (h *TaskHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Error("Failed to encode JSON", "error", err)
	}
}

// CreateTask обрабатывает POST /tasks.
// Создаёт новую задачу и возвращает её ID.
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling POST /tasks")
	id, err := h.service.CreateTask(r.Context())
	if err != nil {
		h.logger.Error("Failed to create task", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// GetTask обрабатывает GET /tasks/{id}.
// Возвращает задачу по ID.
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.logger.Error("Missing task ID")
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing task ID"})
		return
	}
	h.logger.Info("Handling GET /tasks/{id}", "task_id", id)

	task, err := h.service.GetTask(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.logger.Warn("Task not found", "task_id", id, "error", err)
			h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
		} else {
			h.logger.Error("Failed to get task", "task_id", id, "error", err)
			h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		}
		return
	}

	response := struct {
		ID          string           `json:"id"`
		Status      model.TaskStatus `json:"status"`
		CreatedAt   string           `json:"created_at"`
		Duration    string           `json:"duration,omitempty"`
		Result      string           `json:"result,omitempty"`
		CompletedAt string           `json:"completed_at,omitempty"`
	}{
		ID:        task.ID,
		Status:    task.Status,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}
	if task.Status == model.StatusCompleted || task.Status == model.StatusFailed {
		if task.CompletedAt != nil {
			response.Duration = task.CompletedAt.Sub(task.CreatedAt).String()
			response.CompletedAt = task.CompletedAt.Format(time.RFC3339)
		}
		response.Result = task.Result
	}

	h.writeJSON(w, http.StatusOK, response)
}

// DeleteTask обрабатывает DELETE /tasks/{id}.
// Удаляет задачу по ID.
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.logger.Error("Missing task ID")
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing task ID"})
		return
	}
	h.logger.Info("Handling DELETE /tasks/{id}", "task_id", id)

	if err := h.service.DeleteTask(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.logger.Warn("Task not found", "task_id", id, "error", err)
			h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
		} else {
			h.logger.Error("Failed to delete task", "task_id", id, "error", err)
			h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTasks обрабатывает GET /tasks.
// Возвращает список всех задач.
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling GET /tasks")
	tasks, err := h.service.ListTasks()
	if err != nil {
		h.logger.Error("Failed to list tasks", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	response := make([]struct {
		ID          string           `json:"id"`
		Status      model.TaskStatus `json:"status"`
		CreatedAt   string           `json:"created_at"`
		Duration    string           `json:"duration,omitempty"`
		Result      string           `json:"result,omitempty"`
		CompletedAt string           `json:"completed_at,omitempty"`
	}, len(tasks))
	for i, task := range tasks {
		response[i] = struct {
			ID          string           `json:"id"`
			Status      model.TaskStatus `json:"status"`
			CreatedAt   string           `json:"created_at"`
			Duration    string           `json:"duration,omitempty"`
			Result      string           `json:"result,omitempty"`
			CompletedAt string           `json:"completed_at,omitempty"`
		}{
			ID:        task.ID,
			Status:    task.Status,
			CreatedAt: task.CreatedAt.Format(time.RFC3339),
		}
		if task.Status == model.StatusCompleted || task.Status == model.StatusFailed {
			if task.CompletedAt != nil {
				response[i].Duration = task.CompletedAt.Sub(task.CreatedAt).String()
				response[i].CompletedAt = task.CompletedAt.Format(time.RFC3339)
			}
			response[i].Result = task.Result
		}
	}

	h.writeJSON(w, http.StatusOK, response)
}
