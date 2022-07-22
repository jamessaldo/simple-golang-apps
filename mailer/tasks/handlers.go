package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"bytes"
	"text/template"

	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
)

type BodylinkEmail struct {
	NAME string
}

func SendEmail(to string, subject string, data interface{}, templateFile string) error {
	result, _ := ParseTemplate(fmt.Sprintf("%s%s", "./templates/", templateFile), data)
	m := gomail.NewMessage()
	m.SetHeader("From", "rizkysr19@zohomail.com")
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "<RECIPIENT CC>", "<RECIPIENT CC NAME>")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)
	// m.Attach(templateFile) // attach whatever you want
	senderPort := 587
	d := gomail.NewDialer("smtp.zoho.com", senderPort, "rizkysr19@zohomail.com", "WAp2T7sYRDzk")
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
	return err
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Println(err)
		return "", err
	}
	return buf.String(), nil
}

func SendEmailVerification(to string, data interface{}) {
	var err error
	template := "email_template_verifikasi.html"
	subject := "sample email"
	err = SendEmail(to, subject, data, template)
	if err == nil {
		fmt.Println("send email '" + subject + "' success")
	} else {
		fmt.Println(err)
	}
}

// HandleWelcomeEmailTask handler for welcome email task.
func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	// Get user ID from given task.
	var data map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return err
	}

	templateData := BodylinkEmail{
		NAME: data["name"].(string),
	}
	to := "jamessaldo19@gmail.com"
	go SendEmailVerification(to, templateData)

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
