package mail

import (
	"bytes"
	"html/template"
	"log"

	"github.com/fadilmartias/firavel/config"
	"gopkg.in/gomail.v2"
)

func SendEmailVerificationEmail(to string, token string) error {
	mailConfig := config.LoadMailConfig()
	tmpl, err := template.ParseFiles("templates/email/verification_email.html")
	if err != nil {
		return err
	}

	data := struct {
		VerificationLink string
	}{
		VerificationLink: mailConfig.AppURL + "/verify-email?token=" + token,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Username, mailConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	return nil
}
