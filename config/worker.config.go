package config

import (
	"sync"
	"time"
)

type JobFunc func() error
type Job struct {
	Handler JobFunc
	Config  JobConfig
	Retry   int
}

type ExecuteAt struct {
	Year   int
	Month  time.Month
	Day    int
	Hour   int
	Minute int
	Second int
}

type JobConfig struct {
	Once  bool
	Every time.Duration
	At    *ExecuteAt
}

var (
	RUN_ONCE_CONFIG         = JobConfig{Once: true}
	RUN_EVERY_HOUR_CONFIG   = JobConfig{Every: 1 * time.Hour}
	RUN_EVERY_MINUTE_CONFIG = JobConfig{Every: 1 * time.Minute}
)

type workerPool struct {
	workerNum int
	jobChan   chan Job
	quitChan  chan struct{}
	wg        sync.WaitGroup
}

func NewWorkerPool(workerNum int) *workerPool {
	return &workerPool{
		workerNum: workerNum,
		jobChan:   make(chan Job),
		quitChan:  make(chan struct{}),
	}
}

func (wp *workerPool) Start() {
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.worker()
		}()
	}
}

func (wp *workerPool) worker() {
	for {
		select {
		case job := <-wp.jobChan:
			wp.executeJob(job.Handler, job.Retry, job.Config)
		case <-wp.quitChan:
			return
		}
	}
}

func (wp *workerPool) Stop() {
	close(wp.quitChan)
	wp.wg.Wait()
}

func (wp *workerPool) AddJob(job Job) {
	wp.jobChan <- job
}

func (wp *workerPool) executeJob(job JobFunc, retry int, jobConfig JobConfig) {
	if jobConfig.Once {
		for retry > 0 {
			err := job()
			if err != nil {
				retry--
				continue
			}
			retry = 0
		}
	} else if jobConfig.Every != 0 {
		for {
			for retry > 0 {
				err := job()
				if err != nil {
					retry--
					continue
				}
				retry = 0
			}
			time.Sleep(jobConfig.Every)
		}
	} else if jobConfig.At != nil {
		for {
			now := time.Now()
			if jobConfig.At.Year == 0 {
				jobConfig.At.Year = now.Year()
			}
			if jobConfig.At.Month == 0 {
				jobConfig.At.Month = now.Month()
			}
			if jobConfig.At.Day == 0 {
				jobConfig.At.Day = now.Day()
			}
			if jobConfig.At.Hour == 0 {
				jobConfig.At.Hour = now.Hour()
			}
			if jobConfig.At.Minute == 0 {
				jobConfig.At.Minute = now.Minute()
			}
			if jobConfig.At.Second == 0 {
				jobConfig.At.Second = now.Second()
			}
			// Set target time for today
			targetTime := time.Date(
				jobConfig.At.Year,
				jobConfig.At.Month,
				jobConfig.At.Day,
				jobConfig.At.Hour,
				jobConfig.At.Minute,
				jobConfig.At.Second,
				0,
				now.Location(),
			)
			until := time.Until(targetTime)
			time.Sleep(until)
			for retry > 0 {
				err := job()
				if err != nil {
					retry--
					continue
				}
				retry = 0
			}
		}
	}
}
