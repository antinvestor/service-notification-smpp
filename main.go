package main

import (
	"fmt"
	notificationV1 "github.com/antinvestor/service-notification-api"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/template-service/config"
	"github.com/antinvestor/template-service/service"
	"github.com/antinvestor/template-service/service/models"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/antinvestor/apis"

	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"
)

func main() {

	serviceName := "template_service"

	var templateConfig config.TemplateConfig
	err := frame.ConfigProcess("", &templateConfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, srv := frame.NewService(serviceName, frame.Config(&templateConfig))
	defer srv.Stop(ctx)

	logger := srv.L()

	serviceOptions := []frame.Option{frame.Datastore(ctx)}
	if templateConfig.DoDatabaseMigrate() {

		srv.Init(serviceOptions...)

		err = srv.MigrateDatastore(ctx, templateConfig.GetDatabaseMigrationPath(),
			models.Template{})

		if err != nil {
			logger.WithError(err).Fatal("could not migrate successfully")
		}
		return
	}

	err = srv.RegisterForJwt(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("could not register for jwt")
		return
	}

	oauth2ServiceHost := templateConfig.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)

	if templateConfig.Oauth2ServiceAudience != "" {
		audienceList = strings.Split(templateConfig.Oauth2ServiceAudience, ",")
	}

	notificationCli, err := notificationV1.NewNotificationClient(ctx,
		apis.WithEndpoint(templateConfig.NotificationServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(srv.JwtClientID()),
		apis.WithTokenPassword(srv.JwtClientSecret()),
		apis.WithAudiences(audienceList...))
	if err != nil {
		logger.WithError(err).Fatal("could not setup notification client")
	}

	profileCli, err := profileV1.NewProfileClient(ctx,
		apis.WithEndpoint(templateConfig.ProfileServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(srv.JwtClientID()),
		apis.WithTokenPassword(srv.JwtClientSecret()),
		apis.WithAudiences(audienceList...))
	if err != nil {
		logger.WithError(err).Fatal("could not setup profile client")
	}

	partitionCli, err := partitionV1.NewPartitionsClient(
		ctx,
		apis.WithEndpoint(templateConfig.PartitionServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(srv.JwtClientID()),
		apis.WithTokenPassword(srv.JwtClientSecret()),
		apis.WithAudiences(audienceList...))
	if err != nil {
		logger.WithError(err).Fatal("could not setup partition client")
	}

	serviceTranslations := frame.Translations("en")
	serviceOptions = append(serviceOptions, serviceTranslations)

	authServiceHandlers := service.NewAuthRouterV1(srv, &templateConfig, profileCli, partitionCli, notificationCli)

	defaultServer := frame.HttpHandler(authServiceHandlers)
	serviceOptions = append(serviceOptions, defaultServer)

	serviceOptions = append(serviceOptions,
		frame.WithPoolConcurrency(100),
		frame.WithPoolCapacity(500),
	)

	srv.Init(serviceOptions...)

	serverPort := templateConfig.ServerPort
	if serverPort == "" {
		serverPort = "7020"
	}

	logger.WithField("port", serverPort).Info(" initiating server operations")
	err = srv.Run(ctx, fmt.Sprintf(":%v", serverPort))
	if err != nil {
		logger.WithError(err).Error("could not run Server ")
	}
}
