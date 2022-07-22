package main

import (
	"log"
	"runtime"

	"nctwo/mailer/tasks"

	"github.com/hibiken/asynq"
)

const (
	// TypeEmailTask is a name of the task type
	// for sending an email.
	TypeEmailTask = "email:task"

	// TypeDelayedEmail is a name of the task type
	// for sending a delayed email.
	TypeDelayedEmail = "email:delayed"
)

func main() {
	runtime.GOMAXPROCS(2)

	// Create and configuring Redis connection.
	redisConnection := asynq.RedisClientOpt{
		Addr: "localhost:6379", // Redis server address
	}

	// Create and configuring Asynq worker server.
	worker := asynq.NewServer(redisConnection, asynq.Config{
		// Specify how many concurrent workers to use.
		Concurrency: 10,
		// Specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})

	// Create a new task's mux instance.
	mux := asynq.NewServeMux()

	// Define a task handler for the email task.
	mux.HandleFunc(
		TypeEmailTask,         // task type
		tasks.HandleEmailTask, // handler function
	)

	// Define a task handler for the delayed email task.
	mux.HandleFunc(
		TypeDelayedEmail,             // task type
		tasks.HandleDelayedEmailTask, // handler function
	)

	// Run worker server.
	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}
}
