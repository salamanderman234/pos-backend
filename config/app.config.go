package config

import (
	"crypto/rand"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// app
var debug bool
var version string
var applicationKey []byte

// mailer
var mailer *gomail.Dialer

// validator
var vld *validator.Validate

// sanitizer
var sanitizer *bluemonday.Policy

// job
var WorkerPool *workerPool

func StartSetup() {
	// set env
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// application setup
	GenerateApplicationKey()
	debug = viper.GetBool("APP_DEBUG")
	version = viper.GetString("APP_VERSION")
	// sanitizer
	sanitizer = bluemonday.UGCPolicy()
	// validator
	vld = validator.New()
	registerValidation()
	// mailer
	setupMailer()
	// connect db
	connectDB()
	// setup worker
	WorkerPool = NewWorkerPool(10)
	WorkerPool.Start()
}

func GenerateApplicationKey() {
	key := make([]byte, APP_KEY_SIZE)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	applicationKey = key
}

func ApplicationKey() []byte {
	return applicationKey
}

func ApplicationDebugStatus() bool {
	return debug
}

func ApplicationVersion() string {
	return version
}

func Validator() *validator.Validate {
	return vld
}

func Sanitizer() *bluemonday.Policy {
	return sanitizer
}

func Mailer() *gomail.Dialer {
	return mailer
}
