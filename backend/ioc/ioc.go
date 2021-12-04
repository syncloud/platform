package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/syncloud/platform/certificate/certbot"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/redirect"
	"log"
	"time"
)

func Init(configDb string) {

	Singleton(func() *config.UserConfig {
		return config.NewUserConfig(configDb, config.OldConfig)
	})
	Singleton(func() *network.Interface { return network.New() })
	Singleton(func(userConfig *config.UserConfig, iface *network.Interface) *cron.CertificateJob {
		return cron.NewCertificateJob(userConfig, iface)
	})
	Singleton(func(userConfig *config.UserConfig) *cron.PortsJob {
		return cron.NewPortsJob(userConfig)
	})
	Singleton(func(job1 *cron.CertificateJob, job2 *cron.PortsJob, userConfig *config.UserConfig) *cron.Cron {
		return cron.New([]cron.Job{job1, job2}, time.Minute*5, userConfig)
	})
	Singleton(func() *identification.Parser { return identification.New() })
	Singleton(func(userConfig *config.UserConfig, identification *identification.Parser, iface *network.Interface) *redirect.Service {
		return redirect.New(userConfig, identification, iface)
	})
	Singleton(func() *config.SystemConfig {
		return config.NewSystemConfig(config.File)
	})
	Singleton(func(redirectService *redirect.Service, userConfig *config.UserConfig, systemConfig *config.SystemConfig) *certbot.Generator {
		return certbot.New(redirectService, userConfig, systemConfig)
	})
}

func Singleton(resolver interface{}) {
	err := container.Singleton(resolver)
	if err != nil {
		log.Fatalln(err)
	}
}

func Resolve(abstraction interface{}) {
	err := container.Resolve(abstraction)
	if err != nil {
		log.Fatalln(err)
	}
}

func Call(abstraction interface{}) {
	err := container.Call(abstraction)
	if err != nil {
		log.Fatalln(err)
	}
}
