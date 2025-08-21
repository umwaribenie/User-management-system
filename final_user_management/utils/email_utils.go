package utils

import (
	"crypto/tls"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(config EmailConfig, toEmail string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Password Reset OTP")
	m.SetBody("text/html", "Your OTP is: <strong>"+otp+"</strong><br>It will expire in 5 minutes.")

	d := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	if config.SMTPHost == "localhost" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}
