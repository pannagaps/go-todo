package controller

import (
	"fmt"
	"net/http"

	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Added for UUID handling
)

// create a db connection
var dbConn = Connect("tasks")

// dbConn.setTableName("tasks")

// ---- CRUD Logic ----
func GetTasks(c *gin.Context) {
	// Initialize variables
	var tasks []models.Task
	var err error

	// Read the optional `id` query parameter
	id := c.Query("id")
	limit := 50

	if id != "" {
		// Parse the `id` into a UUID
		parsedID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}

		// Fetch the task by ID
		tasks, err = dbConn.Query(limit, parsedID)
		if err != nil {
			if err.Error() == "task not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  fmt.Sprintf("Error fetching task: %v", err),
				})
			}
			return
		}

		// Return the task
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   tasks,
		})
	} else {
		// Fetch all tasks if `id` is not provided
		tasks, err = dbConn.Query(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Error fetching tasks: %v", err),
			})
			return
		}

		// Return the list of tasks
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   tasks,
		})
	}
}

func CreateTask(c *gin.Context) {
	var newTask models.Task
	// Bind JSON from request body to struct
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := dbConn.Insert(newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Error creating tasks: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   res,
	})
}

func UpdateTask(c *gin.Context) {
	var updated models.Task
	id := c.Query("id")
	if id != "" {
		parsedID, err := uuid.Parse(id) // Parse string to UUID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
		// Bind JSON from request body to struct
		if err := c.BindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := dbConn.Update(parsedID, updated)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Error updating tasks: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   res,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is Required"})
		return

	}
	// http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		parsedID, err := uuid.Parse(id) // Parse string to UUID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
		res, err := dbConn.Delete(parsedID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Error deleting tasks: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   res,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is Required"})
		return

	}
}
