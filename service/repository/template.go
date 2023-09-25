package repository

import (
	"context"
	"github.com/antinvestor/template-service/service/models"

	"github.com/pitabwire/frame"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	GetByID(id string) (*models.Template, error)
	GetByName(name string) (*models.Template, error)
	Save(language *models.Template) error
}

type templateRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewTemplateRepository(ctx context.Context, service *frame.Service) TemplateRepository {
	return &templateRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *templateRepository) GetByName(name string) (*models.Template, error) {
	var template models.Template
	err := repo.readDb.First(&template, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (repo *templateRepository) GetByID(id string) (*models.Template, error) {
	template := models.Template{}
	err := repo.readDb.First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (repo *templateRepository) Save(language *models.Template) error {
	return repo.writeDb.Save(language).Error
}
