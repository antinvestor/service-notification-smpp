package events

import (
	"context"
	"errors"

	"github.com/antinvestor/service-notification-smpp/service/models"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/util"
	"gorm.io/gorm/clause"
)

const TemplateSaveEvent = "template.save"

type TemplateSave struct {
	dbPool pool.Pool
}

func NewTemplateSave(_ context.Context, dbPool pool.Pool) *TemplateSave {
	return &TemplateSave{dbPool: dbPool}
}

func (e *TemplateSave) Name() string {
	return TemplateSaveEvent
}

func (e *TemplateSave) PayloadType() any {
	return &models.Template{}
}

func (e *TemplateSave) Validate(ctx context.Context, payload any) error {
	template, ok := payload.(*models.Template)
	if !ok {
		return errors.New(" payload is not of type models.Template")
	}

	if template.GetID() == "" {
		return errors.New(" template Id should already have been set ")
	}

	return nil
}

func (e *TemplateSave) Execute(ctx context.Context, payload any) error {
	template := payload.(*models.Template)

	logger := util.Log(ctx).With("type", e.Name())
	logger.With("payload", template).Info("handling event")

	result := e.dbPool.DB(ctx, false).Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(template)

	err := result.Error
	if err != nil {
		logger.WithError(err).Warn("could not save to db")
		return err
	}
	logger.With("rows affected", result.RowsAffected).Info("successfully saved record to db")

	return nil
}
