package worker

import (
	"encoding/json"
	"time"

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

type Payload struct {
	UserName     string
	TemplateName string
	To           string
}

// NewEmailTask task payload for a new email.
func NewEmailTask(data *Payload) *asynq.Task {
	// Specify task payload.

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeEmailTask, b)
}

// NewDelayedEmailTask task payload for a delayed email.
func NewDelayedEmailTask(id int, ts time.Time) *asynq.Task {
	// Specify task payload.
	payload := map[string]interface{}{
		"user_id": id,          // set user ID
		"sent_in": ts.String(), // set time to sending
	}

	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeDelayedEmail, b)
}
