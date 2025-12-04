package main

import (
	"fmt"
	"os"
	
	"github.com/easyhire/backend/internal/services"
)

func main() {
	fmt.Println("🧪 Прямое тестирование EmailService из Task #9")
	
	// Убедимся что SMTP не настроен (тестовый режим)
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	
	// Создаем email service
	emailService := services.NewEmailService()
	
	// Тестируем отправку
	fmt.Println("\n📧 Тест отправки приглашения:")
	err := emailService.SendInvitation(
		"test.candidate@example.com",
		"inv_test_1234567890", 
		"Senior Go Developer Assessment",
	)
	
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Println("✅ Email service работает корректно!")
		fmt.Println("   В логах должно быть:")
		fmt.Println("   [EMAIL TEST] Invitation for: Senior Go Developer Assessment")
		fmt.Println("     To: test.candidate@example.com")
		fmt.Println("     Token: inv_test_1234567890")
	}
}
