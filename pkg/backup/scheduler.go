package backup

import (
	"github.com/robfig/cron"
)

type Schedule struct {
	Spec string
	Job  *Job
}

type Scheduler struct {
	cr *cron.Cron
}

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

func (s *Scheduler) Run() {
	s.cr.Run()
}
