package db

import (
	"context"
	"html/template"
	"profile_service/pkg/models"
)

type HTMLTemplateDAO interface {
	SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, params []string, template *template.Template) (*models.HTMLTeplate, error)
	GetTemplatesList(ctx context.Context) ([]*models.HTMLTeplate, error)
}
