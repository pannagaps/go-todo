package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"example.com/m/controller"
	"github.com/joho/godotenv"
)


// ---- main entry ----
func main() {

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read values using os.Getenv
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback if not set
	}

	log.SetPrefix("APP: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	fmt.Println("Server running on http://localhost:3000")
	log.Printf("Server running n port 3000")
	http.HandleFunc("/tasks", controller.HandleTasks)
	http.HandleFunc("/tasks/", controller.HandleTaskByID)

	http.ListenAndServe(":"+port, nil)
}
