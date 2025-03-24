package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connection *gorm.DB

func getDsn() string {
	host := viper.GetString("DB_HOST")
	user := viper.GetString("DB_USER")
	pass := viper.GetString("DB_PASS")
	port := viper.GetString("DB_PORT")
	name := viper.GetString("DB_NAME")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)
}

func connectDB() {
	dsn := getDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	connection = db
}

func Conn() *gorm.DB {
	return connection
}

var logConnection *gorm.DB

func getLogDsn() string {
	host := viper.GetString("LOG_DB_HOST")
	user := viper.GetString("LOG_DB_USER")
	pass := viper.GetString("LOG_DB_PASS")
	port := viper.GetString("LOG_DB_PORT")
	name := viper.GetString("LOG_DB_NAME")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)
}

func logConnectDB() {
	dsn := getLogDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	logConnection = db
}

func LogConn() *gorm.DB {
	return logConnection
}
