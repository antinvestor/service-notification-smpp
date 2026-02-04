package main

import (
	"context"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/partition/connectrpc/go/partition/v1/partitionv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	apis "github.com/antinvestor/apis/go/common"
	"github.com/antinvestor/apis/go/notification"
	"github.com/antinvestor/apis/go/partition"
	"github.com/antinvestor/apis/go/profile"
	"github.com/antinvestor/service-notification-smpp/config"
	"github.com/antinvestor/service-notification-smpp/service"
	"github.com/antinvestor/service-notification-smpp/service/events"
	"github.com/antinvestor/service-notification-smpp/service/models"
	"github.com/pitabwire/frame"
	fconfig "github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/openid"
	"github.com/pitabwire/util"
)

func main() {
	tmpCtx := context.Background()

	cfg, err := fconfig.LoadWithOIDC[config.TemplateConfig](tmpCtx)
	if err != nil {
		util.Log(tmpCtx).With("err", err).Error("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = "template_service"
	}

	ctx, svc := frame.NewServiceWithContext(
		tmpCtx,
		frame.WithConfig(&cfg),
		frame.WithRegisterServerOauth2Client(),
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)

	log := util.Log(ctx)
	dbManager := svc.DatastoreManager()

	if cfg.DoDatabaseMigrate() {
		dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)
		if dbPool == nil {
			log.Fatal("database pool is nil - check DATABASE_URL environment variable")
			return
		}
		err = dbManager.Migrate(ctx, dbPool, cfg.GetDatabaseMigrationPath(),
			models.Template{})
		if err != nil {
			log.WithError(err).Fatal("could not migrate successfully")
		}
		return
	}

	sm := svc.SecurityManager()

	audienceList := cfg.GetOauth2ServiceAudience()

	profileCli, err := setupProfileClient(ctx, sm, cfg, audienceList)
	if err != nil {
		log.WithError(err).Fatal("could not setup profile client")
	}

	partitionCli, err := setupPartitionClient(ctx, sm, cfg, audienceList)
	if err != nil {
		log.WithError(err).Fatal("could not setup partition client")
	}

	notificationCli, err := setupNotificationClient(ctx, sm, cfg, audienceList)
	if err != nil {
		log.WithError(err).Fatal("could not setup notification client")
	}

	dbPool := dbManager.GetPool(ctx, datastore.DefaultPoolName)
	if dbPool == nil {
		log.Fatal("database pool is nil - check DATABASE_URL environment variable")
		return
	}

	authServiceHandlers := service.NewAuthRouterV1(svc, &cfg, profileCli, partitionCli, notificationCli)

	serviceOptions := []frame.Option{
		frame.WithHTTPHandler(authServiceHandlers),
		frame.WithRegisterEvents(
			events.NewTemplateSave(ctx, dbPool),
		),
	}

	svc.Init(ctx, serviceOptions...)

	serverPort := cfg.Port()
	if serverPort == "" {
		serverPort = ":7020"
	}

	log.With("port", serverPort).Info("initiating server operations")
	err = svc.Run(ctx, serverPort)
	if err != nil {
		log.WithError(err).Error("could not run Server")
	}
}

func setupProfileClient(
	ctx context.Context,
	clHolder security.InternalOauth2ClientHolder,
	cfg config.TemplateConfig,
	audiences []string,
) (profilev1connect.ProfileServiceClient, error) {
	return profile.NewClient(ctx,
		apis.WithEndpoint(cfg.ProfileServiceURI),
		apis.WithTokenEndpoint(cfg.GetOauth2TokenEndpoint()),
		apis.WithTokenUsername(clHolder.JwtClientID()),
		apis.WithTokenPassword(clHolder.JwtClientSecret()),
		apis.WithScopes(openid.ConstSystemScopeInternal),
		apis.WithAudiences(audiences...))
}

func setupPartitionClient(
	ctx context.Context,
	clHolder security.InternalOauth2ClientHolder,
	cfg config.TemplateConfig,
	audiences []string,
) (partitionv1connect.PartitionServiceClient, error) {
	return partition.NewClient(ctx,
		apis.WithEndpoint(cfg.PartitionServiceURI),
		apis.WithTokenEndpoint(cfg.GetOauth2TokenEndpoint()),
		apis.WithTokenUsername(clHolder.JwtClientID()),
		apis.WithTokenPassword(clHolder.JwtClientSecret()),
		apis.WithScopes(openid.ConstSystemScopeInternal),
		apis.WithAudiences(audiences...))
}

func setupNotificationClient(
	ctx context.Context,
	clHolder security.InternalOauth2ClientHolder,
	cfg config.TemplateConfig,
	audiences []string,
) (notificationv1connect.NotificationServiceClient, error) {
	return notification.NewClient(ctx,
		apis.WithEndpoint(cfg.NotificationServiceURI),
		apis.WithTokenEndpoint(cfg.GetOauth2TokenEndpoint()),
		apis.WithTokenUsername(clHolder.JwtClientID()),
		apis.WithTokenPassword(clHolder.JwtClientSecret()),
		apis.WithScopes(openid.ConstSystemScopeInternal),
		apis.WithAudiences(audiences...))
}
