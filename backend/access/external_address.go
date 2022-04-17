package access

import (
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"go.uber.org/zap"
)

type ExternalAddress struct {
	userConfig *config.UserConfig
	redirect   *redirect.Service
	trigger    *event.Trigger
	logger     *zap.Logger
}

func New(userConfig *config.UserConfig, redirect *redirect.Service, trigger *event.Trigger, logger *zap.Logger) *ExternalAddress {
	return &ExternalAddress{
		userConfig: userConfig,
		redirect:   redirect,
		trigger:    trigger,
		logger:     logger,
	}
}

func (a *ExternalAddress) Update(request model.Access) error {

	a.logger.Info(fmt.Sprintf("set dns, ipv4: %v, ipb4 public: %v, ipv6: %v", request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled))
	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			request.Ipv4,
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
