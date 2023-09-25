package business

import (
	"context"
	partapi "github.com/antinvestor/service-partition-api"
	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/template-service/service/events"
	"github.com/antinvestor/template-service/service/models"
	"github.com/antinvestor/template-service/service/repository"
	"github.com/pitabwire/frame"
)

type TemplateBusiness interface {
	Get(ctx context.Context, id string) (*models.Template, error)
	Store(ctx context.Context, name string) (*models.Template, error)
}

func NewTemplateBusiness(ctx context.Context, service *frame.Service, profileCli *profileV1.ProfileClient, partitionCli *partapi.PartitionClient) (TemplateBusiness, error) {

	if service == nil || profileCli == nil || partitionCli == nil {
		return nil, ErrorInitializationFail
	}

	return &templateBusiness{
		service:    service,
		profileCli: profileCli,
	}, nil
}

type templateBusiness struct {
	service    *frame.Service
	profileCli *profileV1.ProfileClient
}

func (nb *templateBusiness) Get(ctx context.Context, id string) (*models.Template, error) {

	templateRepository := repository.NewTemplateRepository(ctx, nb.service)
	template, err := templateRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return template, nil

}

func (nb *templateBusiness) Store(ctx context.Context, name string) (*models.Template, error) {

	template := models.Template{
		Name: name,
	}

	authClaims := frame.ClaimsFromContext(ctx)
	if authClaims != nil {

		template.TenantID = authClaims.TenantID
		template.PartitionID = authClaims.PartitionID
		template.AccessID = authClaims.AccessID
	}

	template.GenID(ctx)
	// Queue in message for further processing
	event := events.TemplateSave{}
	err := nb.service.Emit(ctx, event.Name(), template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}
