package controller
import (
	"net/http"
	"strconv"
	"strings"
	"example.com/m/utlities"
	
)

// ---- CRUD Logic ----

func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Sort tasks by created time (newest first)
	utlities.WriteJSON(w, http.StatusOK, Tasks)
}
func GetTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, t := range Tasks {
		if t.ID == id {
			utlities.WriteJSON(w, http.StatusOK, t)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}


// Handler for /tasks
func HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTasks(w, r)
	case "POST":
		CreateTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler for /tasks/{id}
func HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		GetTaskByID(w, r, id)
	case "PUT":
		UpdateTask(w, r, id)
	case "DELETE":
		DeleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

