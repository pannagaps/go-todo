// Task model
package models

import "github.com/google/uuid"

type Task struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Duedate     string    `json:"duedate"` // Ensure this matches the GraphQL response field
}
