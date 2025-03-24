package repositories

import (
	"context"

	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/models"
	"gorm.io/gorm"
)

func LogCreate(data models.Log, driver *gorm.DB) error {
	return driver.Create(data).Error
}

func LogRetrieve(ctx context.Context, driver *gorm.DB, logType config.LogTypeEnum, query string, ranges ...int64) ([]models.Log, error) {
	container := []models.Log{}
	lowRange := int64(0)
	highRange := int64(0)

	if len(ranges) > 0 {
		lowRange = ranges[0]
	}

	if len(ranges) > 1 {
		lowRange = ranges[1]
	}

	q := driver.WithContext(ctx).Where("log_type = ?", string(logType)).Where("data LIKE ?", "%"+query+"%")
	if lowRange != 0 {
		q = q.Where("date >= ?", lowRange)
	}
	if highRange != 0 {
		q = q.Where("date <= ?", highRange)
	}
	err := q.Find(&container).Error
	return container, err
}
