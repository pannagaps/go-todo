package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"example.com/m/models"
	"github.com/google/uuid" // Added for UUID handling
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
)

type DbClient struct {
	dbUrl     string
	tableName string
	client    *graphql.Client
}

// --- Singleton Variables ---
var dbInstance *DbClient
var once sync.Once

func Connect(tableName string) *DbClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbUrl := os.Getenv("GRAPHQL_URL")
	client := graphql.NewClient(dbUrl)

	//this guy will run only once
	once.Do(func() {
		fmt.Println("Creating new DBClient instance...")
		dbInstance = &DbClient{
			dbUrl:     dbUrl,
			client:    client,
			tableName: tableName,
		}

	})
	return dbInstance
}

func (c *DbClient) GetTableName() string {
	return c.tableName
}

func (c *DbClient) SetTableName(tableName string) {
	c.tableName = tableName
}

func (c *DbClient) Query(limit int, id ...uuid.UUID) ([]models.Task, error) {
	// Define your GraphQL query
	var query string
	if len(id) > 0 {
		uuid := id[0]
		query = fmt.Sprintf(`
		query MyQuery {
			tasks_by_pk(id: "%s") {
			id	
			title
			description
			completed
			duedate
			}
		}`, uuid)
		// Create a request
		req := graphql.NewRequest(query)
		type singleResponse struct {
			TodoTasksByPk *models.Task `json:"tasks_by_pk"` // Use a pointer to check for nil
		}
		var singleRes singleResponse

		// Execute the query
		if err := c.client.Run(context.Background(), req, &singleRes); err != nil {
			log.Printf("Failed to run query for UUID %s: %v", uuid, err)
			return nil, fmt.Errorf("failed to run query: %w", err)
		}

		// Check if the task is found
		if singleRes.TodoTasksByPk == nil {
			return nil, fmt.Errorf("task not found")
		}

		return []models.Task{*singleRes.TodoTasksByPk}, nil
	} else {
		query = fmt.Sprintf(`
			query MyQuery {
				tasks(limit: %d) {
				id	
				title
				description
				completed
				duedate
				}
			}
		`, limit)

		// Create a request
		req := graphql.NewRequest(query)
		type tasksResponse struct {
			Tasks []models.Task `json:"tasks"`
		}
		var response tasksResponse

		// Execute the query
		if err := c.client.Run(context.Background(), req, &response); err != nil {
			fmt.Printf("Failed to run query with limit %d: %v", limit, err)
			return nil, fmt.Errorf("failed to run query: %w", err)
		}

		return response.Tasks, nil
	}
}

func (c *DbClient) Insert(task models.Task) (response string, error error) {
	// Generate a new UUID for the task
	newUUID := uuid.New()
	task.Id = newUUID

	fmt.Printf("Generated UUID for new task: %s\n", newUUID)

	// Define your GraphQL mutation
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			insert_tasks_one(object: {id: "%s", title: "%s", description: "%s", completed: %t, duedate: "%s"}) {
			id
			}
		}
	`, task.Id, task.Title, task.Description, task.Completed, task.Duedate)

	// Create a request
	req := graphql.NewRequest(mutation)

	var insertresponse struct {
		Insert struct {
			ID uuid.UUID `json:"id"`
		} `json:"insert_tasks_one"`
	}

	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &insertresponse); err != nil {
		log.Printf("Failed to run mutation: %v", err)
		return "", fmt.Errorf("failed to run mutation %w", err)
	}
	return fmt.Sprint("Insert successful:", newUUID), nil
}

func (c *DbClient) Update(id uuid.UUID, task models.Task) (response string, error error) {
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			update_tasks_by_pk(pk_columns: {id: "%s"}, _set: {title: "%s", description: "%s", completed: %t, duedate: "%s"}) {
			id
			}
		}
	`, id, task.Title, task.Description, task.Completed, task.Duedate)

	// Create a request
	req := graphql.NewRequest(mutation)

	var Updateresponse struct {
		Insert struct {
			ID uuid.UUID `json:"id"`
		} `json:"insert_todo_tasks_one"`
	}
	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &Updateresponse); err != nil {
		log.Printf("Failed to update task with UUID %s: %v", id, err)
		return "", fmt.Errorf("failed to run mutation %w", err)
	}
	return fmt.Sprint("Update successful:", id), nil
}

func (c *DbClient) Delete(id uuid.UUID) (response string, error error) {
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			delete_tasks_by_pk(id: "%s") {
			id
			}
		}
	`, id)

	// Create a request
	req := graphql.NewRequest(mutation)

	fmt.Println(mutation)
	var deleteresponse struct {
		Insert struct {
			ID uuid.UUID `json:"id"`
		} `json:"delete_todo_tasks_one"`
	}
	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &deleteresponse); err != nil {
		log.Printf("Failed to delete task with UUID %s: %v", id, err)
		return "", fmt.Errorf("failed to delete:%w", err)
	}
	return fmt.Sprint("Delete successful:", id), nil
}
