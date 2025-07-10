// longtask.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TaskStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var (
	taskStatusMap = make(map[string]string)
	mu            sync.Mutex
)

// StartLongRunningTask starts a long-running task in the background
func StartLongRunningTask(w http.ResponseWriter, r *http.Request) {
	taskID := fmt.Sprintf("%d", time.Now().UnixNano()) // Unique task ID
	mu.Lock()
	taskStatusMap[taskID] = "In Progress"
	mu.Unlock()

	go func(id string) {
		// Simulate a long-running task
		time.Sleep(10 * time.Second) // Simulate work

		mu.Lock()
		taskStatusMap[id] = "Completed"
		mu.Unlock()
	}(taskID)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(TaskStatus{ID: taskID, Status: "Task started"})
}

// CheckTaskStatus checks the status of a task
func CheckTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	mu.Lock()
	status, exists := taskStatusMap[taskID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(TaskStatus{ID: taskID, Status: status})
}
