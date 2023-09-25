package business_test

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/template-service/config"
	"github.com/antinvestor/template-service/service/business"
	"github.com/antinvestor/template-service/service/events"
	"github.com/antinvestor/template-service/service/models"
	"github.com/antinvestor/template-service/service/repository"
	"github.com/golang/mock/gomock"
	"github.com/pitabwire/frame"
	"testing"
)

func getService(serviceName string) *ctxSrv {

	dbURL := frame.GetEnv("TEST_DATABASE_URL", "postgres://ant:secret@localhost:5431/service_template?sslmode=disable")
	testDb := frame.DatastoreCon(dbURL, false)

	var ncfg config.TemplateConfig
	_ = frame.ConfigProcess("", &ncfg)

	ctx, service := frame.NewService(serviceName, testDb, frame.Config(&ncfg), frame.NoopDriver())

	eventList := frame.RegisterEvents(
		&events.TemplateSave{Service: service},
	)
	service.Init(eventList)
	_ = service.Run(ctx, "")
	return &ctxSrv{
		ctx, service,
	}
}

type ctxSrv struct {
	ctx context.Context
	srv *frame.Service
}

func getProfileCli(t *testing.T) *profileV1.ProfileClient {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProfileService := profileV1.NewMockProfileServiceClient(ctrl)
	profileCli := profileV1.InstantiateProfileClient(nil, mockProfileService)
	return profileCli
}

func getPartitionCli(t *testing.T) *partitionV1.PartitionClient {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPartitionService := partitionV1.NewMockPartitionServiceClient(ctrl)

	mockPartitionService.EXPECT().
		GetAccess(gomock.Any(), gomock.Any()).
		Return(&partitionV1.AccessObject{
			AccessId: "test_access-id",
			Partition: &partitionV1.PartitionObject{
				PartitionId: "test_partition-id",
				TenantId:    "test_tenant-id",
			},
		}, nil).AnyTimes()

	profileCli := partitionV1.InstantiatePartitionsClient(nil, mockPartitionService)
	return profileCli
}

func TestNewTemplateBusiness(t *testing.T) {

	profileCli := getProfileCli(t)
	partitionCli := getPartitionCli(t)

	type args struct {
		ctxService   *ctxSrv
		profileCli   *profileV1.ProfileClient
		partitionCli *partitionV1.PartitionClient
	}
	tests := []struct {
		name      string
		args      args
		want      business.TemplateBusiness
		expectErr bool
	}{

		{name: "NewTemplateBusiness",
			args: args{
				ctxService:   getService("NewTemplateBusinessTest"),
				profileCli:   profileCli,
				partitionCli: partitionCli},
			expectErr: false},

		{name: "NewTemplateBusinessWithNils",
			args: args{
				ctxService: &ctxSrv{
					ctx: context.Background(),
				},
				profileCli: nil,
			},
			expectErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := business.NewTemplateBusiness(tt.args.ctxService.ctx, tt.args.ctxService.srv, tt.args.profileCli, tt.args.partitionCli); !tt.expectErr && (err != nil || got == nil) {
				t.Errorf("NewNotificationBusiness() = could not get a valid notificationBusiness at %s", tt.name)
			}
		})
	}
}

func Test_notificationBusiness_Get(t *testing.T) {

	type fields struct {
		ctxService  *ctxSrv
		profileCli  *profileV1.ProfileClient
		partitionCl *partitionV1.PartitionClient
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "NormalGetIDId",
			fields: fields{
				ctxService:  getService("NormalGetIdTest"),
				profileCli:  getProfileCli(t),
				partitionCl: getPartitionCli(t),
			},
			args: args{
				message: "123456",
			},
			wantErr: false,
			want:    "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpl := &models.Template{
				BaseModel: frame.BaseModel{
					ID: tt.args.message,
				},
				Name: "TestingMate",
			}

			trepo := repository.NewTemplateRepository(tt.fields.ctxService.ctx, tt.fields.ctxService.srv)
			err := trepo.Save(tmpl)
			if err != nil {
				t.Errorf("Get() error = %v saving model", err)
				return
			}

			template, err := business.NewTemplateBusiness(
				tt.fields.ctxService.ctx, tt.fields.ctxService.srv, tt.fields.profileCli, tt.fields.partitionCl)
			if err != nil {
				t.Errorf("Get() we could not initiate a new business")
			}
			got, err := template.Get(tt.fields.ctxService.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.GetID() != tt.want {
				t.Errorf("Get() expecting id %s to be reused, got : %s", tt.want, got.GetID())
			}
		})
	}
}
