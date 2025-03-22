package jobs

import (
	"time"
)



type ScheduleJob struct {
	Every 		time.Duration
	EverydayAt 	time.Time
	EveryDate	time.Time
	Job		func() error
}

func EveryMinute(job func() error, minute ...int) ScheduleJob {
	t := 1
	if len(minute) >= 1 {
		t = minute[0]
	}
	return ScheduleJob{
		Every:	time.Duration(t) * time.Minute,
		Job:	job,
	}
}


