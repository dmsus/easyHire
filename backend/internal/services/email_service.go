package services

import (
	"fmt"
	"log"
	"os"
)

type EmailService struct {
	enabled bool
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USERNAME")
	
	enabled := host != "" && port != "" && user != ""
	
	if enabled {
		log.Println("✅ Email service: ENABLED (SMTP configured)")
	} else {
		log.Println("📧 Email service: DISABLED - will log to console")
	}
	
	return &EmailService{enabled: enabled}
}

func (s *EmailService) SendInvitation(email, token, assessmentTitle string) error {
	if !s.enabled {
		// Тестовый режим
		fmt.Printf("[EMAIL TEST] Invitation for: %s\n", assessmentTitle)
		fmt.Printf("  To: %s\n", email)
		fmt.Printf("  Token: %s\n", token)
		fmt.Printf("  Link: http://localhost:3000/assessment/%s\n\n", token)
		return nil
	}
	
	// TODO: Реальная отправка через SMTP
	log.Printf("[EMAIL] Would send real email to: %s", email)
	return nil
}
