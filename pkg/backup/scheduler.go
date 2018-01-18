package backup

import (
	"github.com/robfig/cron"
)

// Schedule represents the schedule of a backup job
type Schedule struct {
	Spec string
	Job  *Job
}

// Scheduler schedules backup jobs
type Scheduler struct {
	cr *cron.Cron
}

// NewScheduler creates a new scheduler or returns an error if the scheduler
// are not valid
func NewScheduler(schedules []Schedule) (*Scheduler, error) {
	cr := cron.New()

	for _, s := range schedules {
		err := cr.AddJob(s.Spec, s.Job)
		if err != nil {
			return nil, err
		}
	}

	return &Scheduler{
		cr: cr,
	}, nil
}

// Run runs the given scheduler, i.e. executing the jobs at their schedule
func (s *Scheduler) Run() {
	s.cr.Run()
}
