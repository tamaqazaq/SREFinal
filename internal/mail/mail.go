package mail

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"os"
)

func SendEmail(to interface{}, subject, body string) error {
	email, ok := to.(string)
	if !ok {
		return fmt.Errorf("invalid email type: expected string, got %T", to)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully to:", email)
	return nil
}
