package config

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func setupMailer() {
	host := viper.GetString("MAIL_HOST")
	port := viper.GetInt("MAIL_PORT")
	user := viper.GetString("MAIL_USER")
	pass := viper.GetString("MAIL_PASS")
	dialer := gomail.NewDialer(host, port, user, pass)
	mailer = dialer
}
