package db

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"profile_service/pkg/models"

	"github.com/google/uuid"
)

var HTMLTemplates map[uuid.UUID]*models.HTMLTeplate

func init() {
	firstUuid := uuid.New()
	template, _ := template.New("test").
		Parse(`<h1>{{ .name }} {{ .age }}<h2>{{ .key}}</h2></h1>`)
	HTMLTemplates = map[uuid.UUID]*models.HTMLTeplate{
		firstUuid: {
			Uuid:     firstUuid,
			Name:     "test",
			Template: template,
		},
	}
}

type HTMLTemplateDAO interface {
	SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, params []string, template *template.Template) (*models.HTMLTeplate, error)
}

type InMemoryTemplateDAO struct {
}

func NewInMemoryTemplateDAO() *InMemoryTemplateDAO {
	return &InMemoryTemplateDAO{}
}

func (dao *InMemoryTemplateDAO) SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, template *template.Template) *models.HTMLTeplate {
	templateObj := &models.HTMLTeplate{
		Uuid:     uuid.New(),
		Name:     templateData.Name,
		Template: template,
	}
	HTMLTemplates[templateObj.Uuid] = templateObj
	return templateObj
}

type FSTemplateDAO struct {
	RootDirectory string
}

func NewFSTemplateDAO(rootDirectory string) *FSTemplateDAO {
	if err := os.MkdirAll(rootDirectory, os.ModePerm); err != nil {
		panic(err)
	}
	return &FSTemplateDAO{RootDirectory: rootDirectory}
}

func (dao *FSTemplateDAO) SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, params []string, template *template.Template) (*models.HTMLTeplate, error) {
	templateObj := &models.HTMLTeplate{
		Uuid:     uuid.New(),
		Name:     templateData.Name,
		Template: template,
	}
	HTMLTemplates[templateObj.Uuid] = templateObj

	file, fileCreateError := os.Create(fmt.Sprintf("%s/%s.json", dao.RootDirectory, templateObj.Uuid.String()))
	if fileCreateError != nil {
		return nil, fileCreateError
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	writeFileError := encoder.Encode(models.HTMLTeplateDumpSchema{
		Uuid:     templateObj.Uuid,
		Name:     templateObj.Name,
		Template: templateData.Template,
		Params:   params,
	})
	// _, writeFileError := w.WriteString("buffered\n")
	if writeFileError != nil {
		return nil, writeFileError
	}
	w.Flush()

	return templateObj, nil
}
