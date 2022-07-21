package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// HandleWelcomeEmailTask handler for welcome email task.
func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	// Get user ID from given task.
	var data map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return err
	}

	// Dummy message to the worker's output.
	fmt.Printf("Send Welcome Email to User ID %d\n", data["user_id"])

	return nil
}

// HandleReminderEmailTask for reminder email task.
func HandleReminderEmailTask(c context.Context, t *asynq.Task) error {
	var data map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return err
	}

	// Dummy message to the worker's output.
	fmt.Printf("Send Reminder Email to User ID %d\n", data["user_id"])
	fmt.Printf("Reason: time is up (%v)\n", data["sent_in"])

	return nil
}
