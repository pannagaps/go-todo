package controller

import (
	"fmt"
	"net/http"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"strconv"

)

// create a db connection
var dbConn = ConnectDb()
// dbConn.setTableName("tasks")

// ---- CRUD Logic ----
func GetTasks(c *gin.Context) {
	dbConn.setTableName("tasks")

	// Sort tasks by created time (newest first)
	var Tasks []models.Task
	var err error

	id:=c.Query("id")
	limit:=50
	if id != "" {
		num, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Fetch by ID
		Tasks, err = dbConn.queryDb(limit, num)
		if err != nil {
			fmt.Printf("Error fetching tasks: %v", err)
			c.JSON(http.StatusInternalServerError, "Failed to fetch tasks")
			return
		}
		c.JSON(http.StatusOK,Tasks)
	} else {
		// Fetch without ID (e.g., get all)
		Tasks, err = dbConn.queryDb(limit) // Or modify queryDb to allow nil / no id
		if err != nil {
			fmt.Printf("Error fetching tasks: %v", err)
			c.JSON(http.StatusInternalServerError, "Failed to fetch tasks")
			return
		}
		c.JSON(http.StatusOK,Tasks)
	}
	
}

func CreateTask(c *gin.Context) {
	dbConn.setTableName("tasks")
	var newTask models.Task
	// Bind JSON from request body to struct
		if err := c.BindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	res, err := dbConn.insertDb(newTask)
	if err != nil {
		fmt.Printf("Error creating tasks: %v", err)
		c.JSON(http.StatusInternalServerError,"Error Creating tasks")
		return
	}
	c.JSON(http.StatusOK,res)
	}

func UpdateTask(c *gin.Context) {
	dbConn.setTableName("tasks")
	var updated models.Task
	id :=c.Query("id")
	if id!=""{
			num, err := strconv.Atoi(id) // Convert string to integer
			if err!=nil{
				fmt.Print("error")
				c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
				return
			}
			// Bind JSON from request body to struct
			if err := c.BindJSON(&updated); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			res, err := dbConn.updateDb(num, updated)
			if err != nil {
				fmt.Printf("Error updating tasks: %v", err)
				c.JSON(http.StatusInternalServerError,"Error Updating Task")
				return
			}
			c.JSON(http.StatusOK,res)
			return
	}else{
			c.JSON(http.StatusBadRequest,gin.H{"error":"id is Required"})
			return

	}
	// http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(c *gin.Context) {
	dbConn.setTableName("tasks")
	id :=c.Query("id")
	if id!=""{
			num, err := strconv.Atoi(id) // Convert string to integer
			if err!=nil{
				fmt.Print("error")
				c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
				return
			}
			res, err := dbConn.deleteDb(num)
			if err != nil {
				fmt.Printf("Error deleting tasks: %v", err)
				c.JSON(http.StatusInternalServerError,"Error Updating Task")
				return
			}
			c.JSON(http.StatusOK,res)
			return
	}else{
			c.JSON(http.StatusBadRequest,gin.H{"error":"id is Required"})
			return

	}
}
