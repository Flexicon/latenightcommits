package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// Scheduler handles running jobs on a given cron Schedule.
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler returns a new Scheduler instance.
func NewScheduler() *Scheduler {
	return &Scheduler{cron.New()}
}

// Run the cron scheduler, or no-op if already running.
func (s *Scheduler) Run() {
	s.cron.Run()
}

// JobDefinition for working with the Scheduler.
type JobDefinition struct {
	Name     string
	Schedule string
	Run      func() error
}

// AddJob adds a job to the Scheduler to be run on the given spec schedule and name.
func (s *Scheduler) AddJob(job *JobDefinition) error {
	if _, err := s.cron.AddFunc(job.Schedule, func() {
		if err := job.Run(); err != nil {
			log.Printf("%s job error: %v", job.Name, err)
		}
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to schedule %s worker", job.Name))
	}
	return nil
}
