package controller

import (
	"context"
	"fmt"
	"os"

	"example.com/m/models"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"

	"sync" //implemented singlton
)

type DbConnector struct {
	dbUrl     string
	tableName string
	client    *graphql.Client
}

// --- Singleton Variables ---
var dbInstance *DbConnector
var once sync.Once

func ConnectDb() *DbConnector {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
	}
	dbUrl := os.Getenv("GRAPHQL_URL")
	client := graphql.NewClient(dbUrl)

	//this guy will run only once
	once.Do(func() {
		fmt.Println("Creating new DBConnector instance...")
		dbInstance = &DbConnector{
			dbUrl:  dbUrl,
			client: client,
		}

	})
	return dbInstance
}

func (c *DbConnector) getTableName() string {
	return c.tableName
}

func (c *DbConnector) setTableName(tableName string) {
	c.tableName = tableName
}

func (c *DbConnector) queryDb(limit int, id ...int) ([]models.Task, error) {
	// Define your GraphQL query
	var query string
	if id != nil {
		num := id[0]
		query = fmt.Sprintf(`
		query MyQuery {
			todo_%s_by_pk(id: %d) {
			id	
			title
			description
			completed
			duedate
			}
		}
	`, c.tableName, num)
		// Create a request
		req := graphql.NewRequest(query)
		// type singleResponse struct {
		// 	Data struct {
		// 		Task models.Task `json:"todo_tasks_by_pk"`
		// 	} `json:"data"`
		// }
		type singleResponse struct {
			TodoTasksByPk models.Task `json:"todo_tasks_by_pk"`
		}
		var singleRes singleResponse

		// Execute the query
		if err := c.client.Run(context.Background(), req, &singleRes); err != nil {
			return nil, fmt.Errorf("failed to run query: %w", err)
		}

		return []models.Task{singleRes.TodoTasksByPk}, nil
	} else {
		query = fmt.Sprintf(`
			query MyQuery {
				todo_%s(limit: %d) {
				id	
				title
				description
				completed
				duedate
				}
			}
		`, c.tableName, limit)

		// Create a request
		req := graphql.NewRequest(query)
		var response map[string][]models.Task

		// Execute the query
		if err := c.client.Run(context.Background(), req, &response); err != nil {
			return nil, fmt.Errorf("failed to run query: %w", err)
		}
		keys := fmt.Sprintf("todo_%s", c.tableName)
		return response[keys], nil
	}
}

func (c *DbConnector) insertDb(task models.Task) (response string, error error) {
	fmt.Println(task)
	// Define your GraphQL mutation
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			insert_todo_%s_one(object: {title: "%s", description: "%s", completed: %t, duedate: "%s"}) {
			id
			}
		}
	`, c.tableName, task.Title, task.Description, task.Completed, task.Duedate)

	// Create a request
	req := graphql.NewRequest(mutation)

	// var response map[string]map[string]models.Task
	var insertresponse struct {
		Insert struct {
			ID int `json:"id"`
		} `json:"insert_todo_tasks_one"`
	}

	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &insertresponse); err != nil {
		return "", fmt.Errorf("failed to run mutation %w", err)
	}
	return "Insert successful", nil
}

func (c *DbConnector) updateDb(id int, task models.Task) (response string, error error) {
	// Define your GraphQL mutation
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			update_todo_%s_by_pk(pk_columns: {id: %d}, _set: {title: "%s", description: "%s", completed: %t, duedate: "%s"}) {
			id
			}
		}
	`, c.tableName, id, task.Title, task.Description, task.Completed, task.Duedate)

	// Create a request
	req := graphql.NewRequest(mutation)

	var Updateresponse struct {
		Insert struct {
			ID int `json:"id"`
		} `json:"insert_todo_tasks_one"`
	}
	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &Updateresponse); err != nil {
		return "", fmt.Errorf("failed to run mutation %w", err)
	}
	return "Update successful", nil
}

func (c *DbConnector) deleteDb(id int) (response string, error error) {
	// Define your GraphQL mutation
	mutation := fmt.Sprintf(`
		mutation MyMutation {
			delete_todo_%s_by_pk(id: %d) {
			id
			}
		}
	`, c.tableName, id)

	// Create a request
	req := graphql.NewRequest(mutation)

	fmt.Println(mutation)
	var deleteresponse struct {
		Insert struct {
			ID int `json:"id"`
		} `json:"delete_todo_tasks_one"`
	}
	// Execute the mutation
	if err := c.client.Run(context.Background(), req, &deleteresponse); err != nil {
		fmt.Printf("failed to delete:%w", err)
		return "", fmt.Errorf("failed to delete:%w", err)

	}
	return "Delete successful", nil

}
