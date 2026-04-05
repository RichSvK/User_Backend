package service

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type smtpConfig struct {
	User    string
	Pass    string
	Host    string
	Port    string
	AppHost string
	AppPort string
	Secret  string
	IsSend  bool
}

func LoadSMTPConfig() (smtpConfig, error) {
	cfg := smtpConfig{
		User:    os.Getenv("SMTP_EMAIL"),
		Pass:    os.Getenv("SMTP_PASSWORD"),
		Host:    os.Getenv("SMTP_HOST"),
		Port:    os.Getenv("SMTP_PORT"),
		AppHost: os.Getenv("APP_HOST"),
		AppPort: os.Getenv("APP_PORT"),
		Secret:  os.Getenv("EMAIL_SECRET_KEY"),
		IsSend:  false,
	}

	if cfg.User == "" || cfg.Pass == "" || cfg.Host == "" ||
		cfg.Port == "" || cfg.AppHost == "" || cfg.AppPort == "" {
		return cfg, errors.New("missing smtp configuration")
	}

	return cfg, nil
}

//go:embed templates/*.html
var templateFS embed.FS

type verificationData struct {
	VerifyURL string
}

func renderVerificationEmail(verifyURL string) (string, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/verification.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, verificationData{VerifyURL: verifyURL}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (cfg *smtpConfig) sendHTML(ctx context.Context, to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)

	msg := []byte(
		"From: " + cfg.User + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
			htmlBody,
	)

	if !cfg.IsSend {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.User, []string{to}, msg)
	}()

	select {
	case <-ctx.Done():
		log.Println("Email send cancelled")
		return ctx.Err()
	case err := <-done:
		log.Println("Email send finished")
		return err
	}
}
