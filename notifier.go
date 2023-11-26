package main

import (
	"log"

	"github.com/spf13/viper"
)

// Notifier encapsulates the logic and infrastructure necessary to send
// notifications to external systems.
type Notifier interface {
	Notify(msg string) error
}

// NoopNotifier simply performs a no-op instead of notifying any external system.
type NoopNotifier struct{}

// Notify returns nil immediately.
func (n *NoopNotifier) Notify(msg string) error {
	return nil
}

// NewNotifier builds a notifier based on the current configuration, or a NoopNotifier
// if none configured.
func NewNotifier() Notifier {
	token := viper.GetString("telegram.bot_token")
	chatID := viper.GetString("telegram.chat_id")

	if token == "" || chatID == "" {
		log.Println("🟨 No notifier configured, using NoopNotifier.")
		return &NoopNotifier{}
	}

	return NewTelegramNotifier(token, chatID)
}
