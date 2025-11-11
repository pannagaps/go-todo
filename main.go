package main

import (
	"fmt"
	"os"
	"example.com/m/controller"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ---- main entry ----
func main() {

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Read values using os.Getenv
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback if not set
	}

	fmt.Println("Server running on port 3000")

	router := gin.Default()
	router.GET("/tasks", controller.GetTasks)
	router.POST("/tasks/create", controller.CreateTask)
	router.PATCH("/tasks", controller.UpdateTask)
	router.DELETE("/tasks", controller.DeleteTask)

	if err := router.Run(); err != nil {
		fmt.Printf("failed to start server:%v", err)
	}
}
