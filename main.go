package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"practical-assignment/internal/handler"
	"practical-assignment/internal/service"
	"practical-assignment/internal/storage"
)

// loggingMiddleware логирует HTTP-запросы.
func loggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Инициализация хранилища и сервиса
	store := storage.NewInMemoryStorage()
	taskService := service.NewTaskService(store, logger)
	// taskService реализует service.TaskProcessor, передаём напрямую
	taskHandler := handler.NewTaskHandler(taskService, logger)

	// Настройка маршрутизатора
	r := mux.NewRouter()
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Применение middleware
	handler := loggingMiddleware(r, logger)

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	logger.Info("Starting server on :8080")

	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()

	// Ожидание SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}
	logger.Info("Server stopped")
}
