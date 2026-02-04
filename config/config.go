package config

import "github.com/pitabwire/frame/config"

type TemplateConfig struct {
	config.ConfigurationDefault
	ProfileServiceURI      string `envDefault:"127.0.0.1:7005" env:"PROFILE_SERVICE_URI"`
	NotificationServiceURI string `envDefault:"127.0.0.1:7005" env:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI    string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`
}
