package db

import (
	"context"
	"profile_service/pkg/models"
)

type HTMLTemplatesDAO interface {
	CreateTemplate(ctx context.Context) *models.HTMLTeplate
}
