package config

import (
	"github.com/spf13/viper"
)

var WorkerPool workerPool

func StartSetup() {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// connect db
	connectDB()
	// setup worker
	WorkerPool = newWorkerPool(10)
	WorkerPool.Start()
}
