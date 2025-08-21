package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	DailyNotificationTemplate = `📢 <b>LateNightCommits Daily Notifier</b> 📊

Fetched a total of %d commits yesterday.

ℹ️ Check out the <a href="https://latenightcommits.com/api/stats">stats page</a> for more.`
)

func runDailyNotification(db *gorm.DB, notifier Notifier) error {
	log.Println("📬 Running Daily Notification job...")

	var sentToday int64
	if err := db.Raw(
		`SELECT COUNT(id) FROM commits WHERE created_at BETWEEN (CURRENT_DATE - INTERVAL 1 DAY) AND CURRENT_DATE;`,
	).Scan(&sentToday).Error; err != nil {
		return errors.Wrap(err, "failed to retrieve daily fetched amount for notifier")
	}

	msg := fmt.Sprintf(DailyNotificationTemplate, sentToday)

	if err := notifier.Notify(msg); err != nil {
		return errors.Wrap(err, "failed to notify")
	}

	log.Println("Successfully sent daily notification!")
	return nil
}
