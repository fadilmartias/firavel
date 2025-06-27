package mail

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	AppURL   string
}

var Config = func() MailConfig {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Fatalf("MAIL_PORT harus berupa angka: %v", err)
	}

	return MailConfig{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		From:     os.Getenv("MAIL_FROM"),
		AppURL:   os.Getenv("APP_URL"),
	}
}()

func SendResetPasswordEmail(to string, token string) error {
	tmpl, err := template.ParseFiles("../../templates/email/reset_password.html")
	if err != nil {
		return err
	}

	data := struct {
		ResetLink string
	}{
		ResetLink: Config.AppURL + "/reset-password?token=" + token,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", Config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Password Request")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(Config.Host, Config.Port, Config.Username, Config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	return nil
}
