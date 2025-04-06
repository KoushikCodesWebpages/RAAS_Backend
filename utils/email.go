package utils

import (
    "log"

    "gopkg.in/mail.v2"
    "github.com/google/uuid"
    "RAAS/config" // make sure this is correctly imported based on your structure
)

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
    UseTLS   bool
}

func SendEmail(cfg EmailConfig, to, subject, body string) error {
    m := mail.NewMessage()
    m.SetHeader("From", cfg.From)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
    d.TLSConfig = nil // Optional: add custom TLS settings if needed

    return d.DialAndSend(m)
}

func GenerateVerificationToken() string {
    return uuid.New().String()
}

func GetEmailConfig() EmailConfig {
    cfg, err := config.InitConfig()
    if err != nil {
        log.Fatalf("Error initializing config: %v", err)
    }

    return EmailConfig{
        Host:     cfg.EmailHost,
        Port:     cfg.EmailPort,
        Username: cfg.EmailHostUser,
        Password: cfg.EmailHostPassword,
        From:     cfg.DefaultFromEmail,
        UseTLS:   cfg.EmailUseTLS,
    }
}
