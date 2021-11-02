package db

import (
	"context"
	"html/template"
	"profile_service/pkg/conf"
	"profile_service/pkg/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBTemplateDAO struct {
	db                 *mongo.Database
	templateCollection string
}

func NewMongoDBTemplateDAO(ctx context.Context, db *mongo.Database, config *conf.Config) *MongoDBTemplateDAO {
	return &MongoDBTemplateDAO{db: db, templateCollection: "templates"}
}

func (dao *MongoDBTemplateDAO) SaveTemplate(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, params []string, template *template.Template) (*models.HTMLTeplate, error) {
	templateObj := &models.HTMLTeplate{
		UUID:     uuid.New(),
		Name:     templateData.Name,
		Template: templateData.Template,
		Params:   params,
	}

	collection := dao.db.Collection(dao.templateCollection)

	_, err := collection.InsertOne(ctx, templateObj)

	if err != nil {
		return nil, err
	}

	return templateObj, nil
}

func (dao *MongoDBTemplateDAO) GetTemplatesList(ctx context.Context) ([]*models.HTMLTeplate, error) {
	var templates []*models.HTMLTeplate

	collection := dao.db.Collection(dao.templateCollection)

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var t models.HTMLTeplate
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}

	return templates, nil
}
