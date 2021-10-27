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
		UUID:     uuid.New(),
		Name:     templateData.Name,
		Template: template,
	}
	HTMLTemplates[templateObj.UUID] = templateObj
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
		UUID:     uuid.New(),
		Name:     templateData.Name,
		Template: template,
	}
	HTMLTemplates[templateObj.UUID] = templateObj

	file, fileCreateError := os.Create(fmt.Sprintf("%s/%s.json", dao.RootDirectory, templateObj.UUID.String()))
	if fileCreateError != nil {
		return nil, fileCreateError
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	writeFileError := encoder.Encode(models.HTMLTeplateDumpSchema{
		UUID:     templateObj.UUID,
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
