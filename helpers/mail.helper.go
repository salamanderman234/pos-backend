package helpers

import (
	"bytes"
	"html/template"

	"github.com/salamanderman234/pos-backend/config"
	"gopkg.in/gomail.v2"
)

func MailSend(to string, subject string, templateName string, data map[string]any) error {
	// Parse html template
	tmpl, err := template.ParseFiles("templates/" + templateName + ".html")
	if err != nil {
		return err
	}

	// Buffer to hold executed template data
	var body bytes.Buffer

	// Inject dynamic data into the template and execute it
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Set up gomail
	m := gomail.NewMessage()
	m.SetHeader("From", "salamanderman234@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set the mail body as rendered HTML template
	m.SetBody("text/html", body.String())

	// Send email
	if err := config.Mailer().DialAndSend(m); err != nil {
		return err
	}

	return nil
}
