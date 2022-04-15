package access

import (
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/redirect"
)

type ExternalAddress struct {
	userConfig *config.UserConfig
	redirect *redirect.Service
	trigger *event.Trigger
}

func New(userConfig *config.UserConfig, redirect *redirect.Service, trigger *event.Trigger) *ExternalAddress {
	return &ExternalAddress{
		userConfig: userConfig,
		redirect: redirect,
		trigger: trigger,
	}
}

func (a *ExternalAddress) Update(external_access bool, manual_public_ip string, manual_access_port int) {
	self.logger.info('set_access: external_access={0}'.format(external_access))

	if self.user_platform_config.is_redirect_enabled():
	self.redirect_service.sync(manual_public_ip, manual_access_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
		self.user_platform_config.get_domain_update_token(), external_access)

	self.user_platform_config.update_device_access(external_access, manual_public_ip, manual_access_port)
	self.event_trigger.trigger_app_event_domain()
}

func (a *ExternalAddress) Sync() {
	update_token = self.user_platform_config.get_domain_update_token()
	if update_token is
None:
	return

	external_access = self.user_platform_config.get_external_access()
	public_ip = self.user_platform_config.get_public_ip()
	manual_access_port = self.user_platform_config.get_manual_access_port()

	self.redirect_service.sync(public_ip, manual_access_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
		update_token, external_access)
}