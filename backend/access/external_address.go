package access

import (
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/http"
	"github.com/syncloud/platform/rest/model"
	"go.uber.org/zap"
)

type UserConfig interface {
	IsRedirectEnabled() bool
	SetIpv4Enabled(enabled bool)
	SetIpv4Public(enabled bool)
	SetIpv6Enabled(enabled bool)
	SetPublicIp(publicIp *string)
	SetPublicPort(port *int)
	GetPublicIp() *string
	GetPublicPort() *int
	IsIpv6Enabled() bool
	IsIpv4Public() bool
	IsIpv4Enabled() bool
}

type Redirect interface {
	Update(ipv4 *string, ipv6 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error
}

type Trigger interface {
	RunAccessChangeEvent() error
}

type NetworkInfo interface {
	IPv6() *string
}

type Response struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

type ExternalAddress struct {
	probe      *PortProbe
	userConfig UserConfig
	redirect   Redirect
	trigger    Trigger
	client     http.Client
	network    NetworkInfo
	logger     *zap.Logger
}

func New(probe *PortProbe, userConfig UserConfig, redirect Redirect, trigger Trigger, client http.Client, network NetworkInfo, logger *zap.Logger) *ExternalAddress {
	return &ExternalAddress{
		probe:      probe,
		userConfig: userConfig,
		redirect:   redirect,
		trigger:    trigger,
		client:     client,
		network:    network,
		logger:     logger,
	}
}

func (a *ExternalAddress) Update(request model.Access) error {

	a.logger.Info(fmt.Sprintf("update ipv4 enabled: %v, ipv4 public: %v, ipv6 enabled: %v",
		request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled))

	if request.Ipv4Enabled {
		port := config.WebAccessPort
		if request.AccessPort != nil {
			port = *request.AccessPort
		}
		err := a.probe.Probe(request.Ipv4, port)
		if err != nil {
			return err
		}
	}

	ipv6 := a.network.IPv6()
	if request.Ipv6Enabled {
		err := a.probe.Probe(ipv6, config.WebAccessPort)
		if err != nil {
			return err
		}
	}

	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			request.Ipv4,
			ipv6,
			request.AccessPort,
			request.Ipv4Enabled,
			request.Ipv4Public,
			request.Ipv6Enabled)
		if err != nil {
			return err
		}
	}
	a.userConfig.SetIpv4Enabled(request.Ipv4Enabled)
	a.userConfig.SetIpv4Public(request.Ipv4Public)
	a.userConfig.SetPublicIp(request.Ipv4)
	a.userConfig.SetIpv6Enabled(request.Ipv6Enabled)
	a.userConfig.SetPublicPort(request.AccessPort)

	return a.trigger.RunAccessChangeEvent()

}

func (a *ExternalAddress) Sync() error {

	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			a.userConfig.GetPublicIp(),
			a.network.IPv6(),
			a.userConfig.GetPublicPort(),
			a.userConfig.IsIpv4Enabled(),
			a.userConfig.IsIpv4Public(),
			a.userConfig.IsIpv6Enabled())
		if err != nil {
			return err
		}
	}
	return nil
}