package repository

import (
	"context"

	"github.com/antinvestor/service-notification-smpp/service/models"
	"github.com/pitabwire/frame/datastore/pool"
)

type TemplateRepository interface {
	GetByID(ctx context.Context, id string) (*models.Template, error)
	GetByName(ctx context.Context, name string) (*models.Template, error)
	Save(ctx context.Context, language *models.Template) error
}

type templateRepository struct {
	dbPool pool.Pool
}

func NewTemplateRepository(_ context.Context, dbPool pool.Pool) TemplateRepository {
	return &templateRepository{dbPool: dbPool}
}

func (repo *templateRepository) GetByName(ctx context.Context, name string) (*models.Template, error) {
	var template models.Template
	err := repo.dbPool.DB(ctx, true).First(&template, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (repo *templateRepository) GetByID(ctx context.Context, id string) (*models.Template, error) {
	template := models.Template{}
	err := repo.dbPool.DB(ctx, true).First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (repo *templateRepository) Save(ctx context.Context, language *models.Template) error {
	return repo.dbPool.DB(ctx, false).Save(language).Error
}
