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

func (a *ExternalAddress) Update(request model.Access) {

	a.logger.Info(fmt.Sprintf("set dns, ipv4: %v, ipb4 public: %v, ipv6: %v", request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled))
	if a.userConfig.IsRedirectEnabled() {
		a.redirect.Update(request.Ipv4, request.AccessPort, request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled)

	}
	/*
		self.user_platform_config.update_device_access(external_access, manual_public_ip, manual_access_port)
		self.event_trigger.trigger_app_event_domain()
	*/
}

func (a *ExternalAddress) Sync() {
	/*
	   	update_token = self.user_platform_config.get_domain_update_token()
	   	if update_token is
	   None:
	   	return

	   	external_access = self.user_platform_config.get_external_access()
	   	public_ip = self.user_platform_config.get_public_ip()
	   	manual_access_port = self.user_platform_config.get_manual_access_port()

	   	self.redirect_service.sync(public_ip, manual_access_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
	   		update_token, external_access)
	*/
}
