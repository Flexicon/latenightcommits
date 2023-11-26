package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func runFetchWorker(db *gorm.DB, api *GitHubAPI, notifier Notifier) error {
	log.Println("Running fetch_worker")
	start := time.Now()

	// Expose simple http liveness check
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Late Night Commits fetch worker running since: %v", start)
		})

		log.Fatalln(http.ListenAndServe(":"+viper.GetString("port"), nil))
	}()

	// Run fetch job on a schedule
	scheduler := NewScheduler()
	jobs := []*JobDefinition{
		{
			Name:     "daily_notifier",
			Schedule: viper.GetString("daily_notifier.schedule"),
			Run: func() error {
				return runDailyNotification(db, notifier)
			},
		},
		{
			Name:     "fetch_commits",
			Schedule: viper.GetString("fetch_worker.schedule"),
			Run: func() error {
				return runFetchJob(db, api)
			},
		},
	}

	for _, job := range jobs {
		log.Printf(`Scheduling %s worker to run on "%s"`, job.Name, job.Schedule)
		if err := scheduler.AddJob(job); err != nil {
			return err
		}
	}

	scheduler.Run()
	return nil
}
