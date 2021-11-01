package db

import (
	"context"
	"html/template"
	"profile_service/pkg/models"

	"github.com/google/uuid"
)

var HTMLTemplates map[uuid.UUID]*models.HTMLTeplate

func init() {
	firstUUID := uuid.New()
	template, _ := template.New("test").
		Parse(`<h1>{{ .name }} {{ .age }}<h2>{{ .key}}</h2></h1>`)
	HTMLTemplates = map[uuid.UUID]*models.HTMLTeplate{
		firstUUID: {
			UUID:     firstUUID,
			Name:     "test",
			Template: template,
		},
	}
}

type InMemoryTemplateDAO struct {
}

func NewInMemoryTemplateDAO() *InMemoryTemplateDAO {
	return &InMemoryTemplateDAO{}
}

func (dao *InMemoryTemplateDAO) SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, template *template.Template) *models.HTMLTeplate {
	templateObj := &models.HTMLTeplate{
		UUID:     uuid.New(),
		Name:     templateData.Name,
		Template: template,
	}
	HTMLTemplates[templateObj.UUID] = templateObj
	return templateObj
}
