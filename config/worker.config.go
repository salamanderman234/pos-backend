package config

import (
	"sync"
)

type workerPool struct {
	workerNum 	int
	jobChan 	chan func() error
	quitChan 	chan struct{}
	wg		sync.WaitGroup
}

func NewWorkerPool(workerNum int) *workerPool {
	return &workerPool{
		workerNum: 	workerNum,
		jobChan: 	make(chan func() error),
		quitChan: 	make(chan struct{}),
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
				job()
			case <-wp.quitChan:
				return
		}
	}
}

func (wp *workerPool) Stop() {
	close(wp.quitChan)
	wp.wg.Wait()
}

func (wp *workerPool) AddJob(job func() error) {
	wp.jobChan <- job
}

