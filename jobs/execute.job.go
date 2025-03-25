package jobs

import (
	"github.com/salamanderman234/pos-backend/config"
)

var tasks = []config.Job{}

// will executed once application start
func StartJob() {
	for _, task := range tasks {
		config.WorkerPool.AddJob(task)
	}
}
