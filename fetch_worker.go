package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func runFetchWorker(db *gorm.DB, api *GitHubAPI) error {
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
	scheduler := cron.New()

	if _, err := scheduler.AddFunc(viper.GetString("fetch_worker.schedule"), func() {
		if err := runFetchJob(db, api); err != nil {
			log.Printf("Fetch job error: %v", err)
		}
	}); err != nil {
		return errors.Wrap(err, "failed to schedule fetch job worker")
	}

	scheduler.Run()
	return nil
}
