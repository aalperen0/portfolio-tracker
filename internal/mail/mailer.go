package mail

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"gopkg.in/mail.v2"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// / Initialize a new mail.Dialer instance with the given SMTP server settings
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient string, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// / Execute the named template "subject", passing in the dynamic data and storing the
	// / result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// Initalize a new mail instance, set headers recipient, sender, subject.
	// Set body as plainBody to the message
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())

	// Opens a connection to the SMTP server, sends the message and close connection
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
