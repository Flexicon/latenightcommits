package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func runFetchWorker(db *gorm.DB, api *GitHubAPI) error {
	log.Println("Running fetch_worker")
	scheduler := cron.New()

	_, err := scheduler.AddFunc("*/10 * * * *", func() {
		if err := runFetchJob(db, api); err != nil {
			log.Printf("Fetch job error: %v", err)
		}
	})
	if err != nil {
		return errors.Wrap(err, "failed to schedule fetch job worker")
	}

	scheduler.Run()
	return nil
}
