package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	WeeklyNotificationTemplate = `📢 *LateNightCommits Weekly Notifier* 📊

Stats for commits fetched in the last 7 days:

%s
Total commits fetched: %s

ℹ️ Check out the [stats page](https://latenightcommits.com/api/stats) for more\.`
)

func runWeeklyNotification(db *gorm.DB, notifier Notifier) error {
	log.Println("📬 Running Weekly Notification job...")

	var stats struct {
		Daily []struct {
			Date  time.Time
			Count int64
		}
		Total int64
	}
	dbSession := db.Session(&gorm.Session{PrepareStmt: true})

	err := dbSession.
		Raw(`
					SELECT COUNT(id) as 'count', DATE(created_at) as 'date'
					FROM commits
					WHERE created_at >= (CURRENT_DATE - INTERVAL 7 DAY)
					GROUP BY date
					ORDER BY date DESC
					LIMIT 7;
				`).Scan(&stats.Daily).Error
	if err != nil {
		return errors.Wrap(err, "failed to retrieve stats")
	}

	err = dbSession.Raw(`SELECT COUNT(id) FROM commits;`).Scan(&stats.Total).Error
	if err != nil {
		return errors.Wrap(err, "failed to retrieve stats total")
	}

	var dailyMessageBuilder strings.Builder
	for _, s := range stats.Daily {
		dailyMessageBuilder.WriteString(
			fmt.Sprintf("📅 %s: %s\n", s.Date.Format("2006-01-02"), formatNumber(s.Count)),
		)
	}

	msg := fmt.Sprintf(
		WeeklyNotificationTemplate,
		dailyMessageBuilder.String(),
		formatNumber(stats.Total),
	)
	if err := notifier.Notify(msg); err != nil {
		return errors.Wrap(err, "failed to notify")
	}

	log.Println("Successfully sent daily notification!")
	return nil
}
