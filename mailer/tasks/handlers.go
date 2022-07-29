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

type Payload struct {
	UserName     string
	TemplateName string
	To           string
}

func SendEmail(to string, subject string, data interface{}, templateFile string) error {
	result, _ := ParseTemplate(fmt.Sprintf("%s%s", "./templates/", templateFile), data)
	m := gomail.NewMessage()
	m.SetHeader("From", "Tangled Team <rizkysr19@zohomail.com>")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)
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

func SendEmailTask(to, templateName string, data interface{}) {
	var err error
	template := templateName
	subject := "Selamat datang di Tangled"
	err = SendEmail(to, subject, data, template)
	if err == nil {
		fmt.Println("Send email '" + subject + "' to '" + to + "' success")
	} else {
		fmt.Println(err)
	}
}

// HandleEmailTask handler for email task.
func HandleEmailTask(c context.Context, t *asynq.Task) error {
	// Get user ID from given task.
	var data map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return err
	}

	templateData := Payload{
		UserName: data["UserName"].(string),
	}
	to := data["To"].(string)
	fmt.Printf("Sending Email to %s\n", data["UserName"].(string))
	go SendEmailTask(to, data["TemplateName"].(string), templateData)

	return nil
}

// HandleDelayedEmailTask for delayed email task.
func HandleDelayedEmailTask(c context.Context, t *asynq.Task) error {
	var data map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return err
	}

	// Dummy message to the worker's output.
	fmt.Printf("Send Delayed Email to %s\n", data["UserName"].(string))
	fmt.Printf("Reason: time is up (%v)\n", data["sent_in"])

	return nil
}
