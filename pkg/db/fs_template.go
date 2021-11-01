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
