package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"practical-assignment/internal/model"
	"practical-assignment/internal/service"
	"time"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.service.CreateTask(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetTask handles GET /tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, exists := h.service.GetTask(id)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	response := struct {
		ID        string           `json:"id"`
		Status    model.TaskStatus `json:"status"`
		CreatedAt string           `json:"created_at"`
		Duration  string           `json:"duration,omitempty"`
		Result    string           `json:"result,omitempty"`
	}{
		ID:        task.ID,
		Status:    task.Status,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}
	if task.CompletedAt != nil {
		response.Duration = task.CompletedAt.Sub(task.CreatedAt).String()
		response.Result = task.Result
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteTask handles DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !h.service.DeleteTask(id) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
