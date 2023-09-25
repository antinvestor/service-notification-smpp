package events

import (
	"context"
	"errors"
	"github.com/antinvestor/template-service/service/models"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type TemplateSave struct {
	Service *frame.Service
}

func (e *TemplateSave) Name() string {
	return "template.save"
}

func (e *TemplateSave) PayloadType() interface{} {
	return &models.Template{}
}

func (e *TemplateSave) Validate(ctx context.Context, payload interface{}) error {
	template, ok := payload.(*models.Template)
	if !ok {
		return errors.New(" payload is not of type models.Template")
	}

	if template.GetID() == "" {
		return errors.New(" template Id should already have been set ")
	}

	return nil
}

func (e *TemplateSave) Execute(ctx context.Context, payload interface{}) error {
	template := payload.(*models.Template)

	logger := logrus.WithField("type", e.Name())
	logger.WithField("payload", template).Info("handling event")

	result := e.Service.DB(ctx, false).Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(template)

	err := result.Error
	if err != nil {
		logger.WithError(err).Warn("could not save to db")
		return err
	}
	logger.WithField("rows affected", result.RowsAffected).Info("successfully saved record to db")

	return nil
}
