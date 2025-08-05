package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/fadilmartias/firavel/app/services"
	"github.com/fadilmartias/firavel/config"
	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeEmailResetPassword, func(ctx context.Context, t *asynq.Task) error {
		var payload EmailResetPasswordPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		retries, ok := asynq.GetRetryCount(ctx)
		if !ok {
			return fmt.Errorf("failed to get retry count")
		}
		maxRetry, ok := asynq.GetMaxRetry(ctx)
		if !ok {
			return fmt.Errorf("failed to get max retry")
		}
		if retries == maxRetry {
			services.TelegramSendMessage("Worker: Gagal mengirim email reset password: " + string(t.Payload()))
			return fmt.Errorf("max retry reached")
		}
		mailConfig := config.LoadMailConfig()
		tmpl, err := template.ParseFiles("templates/email/reset_password.html")
		if err != nil {
			return err
		}

		data := struct {
			ResetLink string
		}{
			ResetLink: mailConfig.AppURL + "/reset-password?token=" + payload.Token,
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, data); err != nil {
			return err
		}

		m := gomail.NewMessage()
		m.SetHeader("From", mailConfig.From)
		m.SetHeader("To", payload.To)
		m.SetHeader("Subject", "Reset Password Request")
		m.SetBody("text/html", body.String())

		d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Username, mailConfig.Password)

		if err := d.DialAndSend(m); err != nil {
			log.Println("Error sending email:", err)
			return err
		}
		return nil
	})

	mux.HandleFunc(TypeEmailVerification, func(ctx context.Context, t *asynq.Task) error {
		var payload EmailVerificationPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		retries, ok := asynq.GetRetryCount(ctx)
		if !ok {
			return fmt.Errorf("failed to get retry count")
		}
		maxRetry, ok := asynq.GetMaxRetry(ctx)
		if !ok {
			return fmt.Errorf("failed to get max retry")
		}
		if retries == maxRetry {
			services.TelegramSendMessage("Worker: Gagal mengirim email verifikasi: " + string(t.Payload()))
			return fmt.Errorf("max retry reached")
		}
		mailConfig := config.LoadMailConfig()

		// Prepare template
		tmpl, err := template.ParseFiles("templates/email/verification_email.html")
		if err != nil {
			return err
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, struct {
			VerificationLink string
		}{
			VerificationLink: mailConfig.AppURL + "/verify-email?token=" + payload.Token,
		}); err != nil {
			return err
		}

		m := gomail.NewMessage()
		m.SetHeader("From", mailConfig.From)
		m.SetHeader("To", payload.To)
		m.SetHeader("Subject", "Email Verification")
		m.SetBody("text/html", body.String())

		d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Username, mailConfig.Password)

		if err := d.DialAndSend(m); err != nil {
			log.Println("Failed to send email:", err)
			return err
		}

		return nil
	})

	return mux
}
