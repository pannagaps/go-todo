package main

import (
	"os"
	"fmt"
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
	router.GET("/gettask", controller.GetTasks)
	router.POST("/createtask", controller.CreateTask)
	router.PUT("/updatetask", controller.UpdateTask)
	router.POST("/deletetask")

	if err := router.Run(); err != nil {
		fmt.Printf("failed to start server:%v", err)
	}
}
