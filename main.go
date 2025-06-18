package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"practical-assignment/internal/handler"
	"practical-assignment/internal/service"
	"practical-assignment/internal/storage"
)

func main() {
	// Initialize storage and service
	store := storage.NewInMemoryStorage()
	taskService := service.NewTaskService(store)
	taskHandler := handler.NewTaskHandler(taskService)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
