package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/m/models"
	"example.com/m/utlities"
)

var nextID = 1
var Tasks []models.Task


func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask.ID = nextID
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()
	nextID++

	Tasks = append(Tasks, newTask)
	utlities.WriteJSON(w, http.StatusCreated, newTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request, id int) {
	var updated models.Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, t := range Tasks {
		if t.ID == id {
			updated.ID = id
			updated.CreatedAt = t.CreatedAt
			updated.UpdatedAt = time.Now()
			Tasks[i] = updated
			utlities.WriteJSON(w, http.StatusOK, updated)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request, id int) {
	for i, t := range Tasks {
		if t.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
