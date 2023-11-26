package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	DailyNotificationTemplate = `📢 *LateNightCommits Daily Notifier* 📊

Fetched a total of %d commits today\.

ℹ️ Check out the [stats page](https://latenightcommits.com/api/stats) for more\.`
)

func runDailyNotification(db *gorm.DB, notifier Notifier) error {
	log.Println("📬 Running Daily Notification job...")

	var sentToday int64
	if err := db.Raw(
		`SELECT COUNT(id) FROM commits WHERE created_at BETWEEN CURRENT_DATE AND (CURRENT_DATE + INTERVAL 1 DAY);`,
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
