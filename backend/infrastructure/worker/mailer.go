package worker

import (
	"log"

	"github.com/hibiken/asynq"
)

type WorkerInterface interface {
	SendEmail(payload *Payload) error
}

type AsynqClient struct {
	client *asynq.Client
}

var _ WorkerInterface = &AsynqClient{}

func NewWorker(client *asynq.Client) *AsynqClient {
	return &AsynqClient{client: client}
}

//Enqueue task to send email
func (ac *AsynqClient) SendEmail(payload *Payload) error {
	// Define tasks.
	task := NewEmailTask(payload)

	// Process the task immediately in critical queue.
	if _, err := ac.client.Enqueue(
		task,                    // task payload
		asynq.Queue("critical"), // set queue for task
	); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// 	delay := 10 * time.Second

// 	// Define tasks.
// 	task2 := tasks.NewReminderEmailTask(userID, time.Now().Add(delay))

// 	// Process the task 2 minutes later in low queue.
// 	if _, err := client.Enqueue(
// 		task2,                  // task payload
// 		asynq.Queue("low"),     // set queue for task
// 		asynq.ProcessIn(delay), // set time to process task
// 	); err != nil {
// 		log.Fatal(err)
// 	}
