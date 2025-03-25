package services

import (
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
)

func MailSend(target string, template string, subject string, data map[string]any) {
	handler := func() error {
		return helpers.MailSend(target, subject, template, data)
	}
	job := config.Job{
		Handler: handler,
		Config:  config.RUN_ONCE_CONFIG,
		Retry:   config.JOB_SEND_MAIL_RETRY,
	}
	config.WorkerPool.AddJob(job)
}

func MailSendVerify(email string, username string, code string) {
	template := "verify.html"
	subject := "Verifikasi akun anda sekarang !"
	data := map[string]any{
		"username": username,
		"code":     code,
	}
	MailSend(email, template, subject, data)
}

func MailSendResetPassword(email string, username string, code string) {
	template := "reset.html"
	subject := "Reset password akun anda"
	data := map[string]any{
		"username": username,
		"code":     code,
	}
	MailSend(email, template, subject, data)
}

func MailSendTwoFactor(email string, code string) {
	template := "two-factor.html"
	subject := "Verifikasi aktifitas login anda"
	data := map[string]any{
		"code": code,
	}
	MailSend(email, template, subject, data)
}

func MailSendUserWarn(email string, username string, message string) {
	template := "user-warn.html"
	subject := "Aktifitas mencurigakan pada akun anda"
	data := map[string]any{
		"message": message,
	}
	MailSend(email, template, subject, data)
}
