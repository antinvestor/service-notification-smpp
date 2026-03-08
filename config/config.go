package config

import "github.com/pitabwire/frame/config"

type TemplateConfig struct {
	config.ConfigurationDefault
	ProfileServiceURI                        string `envDefault:"127.0.0.1:7005" env:"PROFILE_SERVICE_URI"`
	NotificationServiceURI                   string `envDefault:"127.0.0.1:7005" env:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI                      string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`
	ProfileServiceWorkloadAPITargetPath      string `envDefault:"/ns/profile/sa/service-profile" env:"PROFILE_SERVICE_WORKLOAD_API_TARGET_PATH"`
	PartitionServiceWorkloadAPITargetPath    string `envDefault:"/ns/auth/sa/service-tenancy" env:"PARTITION_SERVICE_WORKLOAD_API_TARGET_PATH"`
	NotificationServiceWorkloadAPITargetPath string `envDefault:"/ns/notifications/sa/service-notification" env:"NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH"`
}
