package repositories

import (
	"context"

	"github.com/salamanderman234/pos-backend/models"
)

func UserFindByUsername(ctx context.Context, username string, selects []string, preloads []string) (models.User, error) {
	return models.User{}, nil
}
func UserFindByID(ctx context.Context, id any, selects []string, preloads []string) (models.User, error) {
	return models.User{}, nil
}

func UserUpdate(ctx context.Context, id any, data models.User, selects []string) (models.User, error) {
	return models.User{}, nil
}

func UserGetLatestPasswordHashs(ctx context.Context, id any) ([]models.UserPasswordHash, error) {
	return []models.UserPasswordHash{}, nil
}
