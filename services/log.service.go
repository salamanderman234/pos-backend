package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/repositories"
)

func LogDispatchLoginAttempt(userID uint, device string, ip string, isSuccess bool, twoFactor bool, token string) {
	data := map[string]any{
		"device":             device,
		"ip":                 ip,
		"is_success":         isSuccess,
		"required_next_step": twoFactor && token == "",
	}
	jsonData, _ := json.Marshal(data)
	message := fmt.Sprintf("Login attempt from %s(%s)", ip, device)
	if twoFactor && token != "" {
		message = fmt.Sprintf("Successfully verify two factor from %s(%s)", ip, device)
	}
	LogDispatch(userID, message, config.LogTypeEnum_FAILURE, string(jsonData))
}

func LogDispatchFailure(userID uint, message string, instance string) {
	data := map[string]any{
		"instance": instance,
		"message":  message,
	}
	jsonData, _ := json.Marshal(data)
	LogDispatch(userID, message, config.LogTypeEnum_FAILURE, string(jsonData))
}

func LogDispatchUserActivity(userID uint, message string) {
	LogDispatch(userID, message, config.LogTypeEnum_USER_ACTIVITY, "")
}

func LogDispatchRequest(userID uint, device string, ip string, method string, url string, status int) {
	data := map[string]any{
		"device": device,
		"ip":     ip,
		"method": method,
		"url":    url,
		"status": status,
	}
	jsonData, _ := json.Marshal(data)
	message := fmt.Sprintf("New request from %s(%s) to (%s)%s", ip, device, method, url)
	LogDispatch(userID, message, config.LogTypeEnum_REQUEST, string(jsonData))
}

func LogDispatchUserChangeLevel(userID uint, old string, new string) {
	data := map[string]any{
		"old": old,
		"new": new,
	}
	jsonData, _ := json.Marshal(data)
	message := fmt.Sprintf("User #%d update their level from %s to %s", userID, old, new)
	LogDispatch(userID, message, config.LogTypeEnum_UPDATE_LEVEL, string(jsonData))
}

func ReadLog(ctx context.Context, logType config.LogTypeEnum, q string, ranges ...int64) ([]models.Log, error) {
	container := []models.Log{}
	err := error(nil)
	driver := config.LogDriver()
	lowRange := int64(0)
	highRange := int64(0)

	if len(ranges) > 0 {
		lowRange = ranges[0]
	}

	if len(ranges) > 1 {
		lowRange = ranges[1]
	}

	switch driver {
	case config.LogDriverEnum_SERVICE:
		urlService := config.LogService()
		params := url.Values{}
		params.Add("type", string(logType))
		params.Add("query", q)
		if lowRange != 0 {
			params.Add("start", fmt.Sprintf("%d", lowRange))
		}
		if highRange != 0 {
			params.Add("end", fmt.Sprintf("%d", highRange))
		}

		fullURL := fmt.Sprintf("%s?%s", urlService, params.Encode())
		resp, err := http.Get(fullURL)
		if err != nil {
			return container, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		result := map[string]any{}
		json.Unmarshal(body, &result)

		list := result["data"]

		listJson, _ := json.Marshal(list)
		json.Unmarshal(listJson, &container)
	case config.LogDriverEnum_EXTERNAL_DATABASE:
		container, err = repositories.LogRetrieve(ctx, config.LogConn(), logType, q, lowRange, highRange)
	default:
		container, err = repositories.LogRetrieve(ctx, config.Conn(), logType, q, lowRange, highRange)
	}
	return container, err
}

func LogDispatch(userID uint, msg string, logType config.LogTypeEnum, data any) {
	job := config.Job{
		Config: config.RUN_ONCE_CONFIG,
		Retry:  config.JOB_LOG_RETRY,
	}
	driver := config.LogDriver()
	now := time.Now().Unix()
	jsonString, _ := json.Marshal(data)
	dbData := models.Log{
		UserID:  userID,
		LogType: string(logType),
		Data:    string(jsonString),
		Message: msg,
		Date:    now,
	}
	switch driver {
	case config.LogDriverEnum_SERVICE:
		job.Handler = func() error {
			url := config.LogService()
			jsonBody, _ := json.Marshal(dbData)
			_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
			return err
		}
	case config.LogDriverEnum_EXTERNAL_DATABASE:
		job.Handler = func() error {
			return repositories.LogCreate(dbData, config.LogConn())
		}
	default:
		job.Handler = func() error {
			return repositories.LogCreate(dbData, config.Conn())
		}
	}
	config.WorkerPool.AddJob(job)
}
