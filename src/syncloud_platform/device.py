import uuid

import requests
from syncloud_platform.config.config import WEB_CERTIFICATE_PORT, WEB_ACCESS_PORT, WEB_PROTOCOL
from syncloudlib import logger, fs

http_network_protocol = 'TCP'


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service,
                 port_drill_factory, platform_cron, ldap_auth, event_trigger, tls, nginx):
        self.tls = tls
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory
        self.auth = ldap_auth
        self.platform_cron = platform_cron
        self.event_trigger = event_trigger
        self.logger = logger.get_logger('Device')
        self.nginx = nginx

    def prepare_redirect(self, redirect_email, redirect_password):

        self.logger.info("prepare redirect {0}".format(redirect_email))
        self.user_platform_config.set_redirect_enabled(True)
        self.user_platform_config.set_user_email(redirect_email)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        return user

    def activate(self, redirect_email, redirect_password, user_domain, device_username, device_password):
        user_domain_lower = user_domain.lower()
        self.logger.info("activate {0}, {1}".format(user_domain_lower, device_username))
        
        self._check_internet_connection()
        user = self.prepare_redirect(redirect_email, redirect_password)
        self.user_platform_config.set_user_update_token(user.update_token)

        main_domain = self.user_platform_config.get_redirect_domain()
        name, email = parse_username(device_username, '{0}.{1}'.format(user_domain_lower, main_domain))
     
        response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain_lower)
        self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)
   
        self._activate_common(name, device_username, device_password, email)

    def activate_custom_domain(self, full_domain, device_username, device_password):
        full_domain_lower = full_domain.lower()
        self.logger.info("activate custom {0}, {1}".format(full_domain_lower, device_username))
        
        self._check_internet_connection()
        
        self.user_platform_config.set_redirect_enabled(False)
        self.user_platform_config.set_custom_domain(full_domain_lower)

        name, email = parse_username(device_username, full_domain_lower)
        self.user_platform_config.set_user_email(email)
       
        self._activate_common(name, device_username, device_password, email)
        
    def _activate_common(self, name, device_username, device_password, email):
    
        self.set_access(False, False, None, 0, 0)

        self.logger.info("activating ldap")
        self.user_platform_config.set_web_secret_key(str(uuid.uuid4().hex))

        self.tls.generate_self_signed_certificate()

        self.auth.reset(name, device_username, device_password, email)
        
        self.nginx.init_config()
        self.nginx.reload_public()

        self.user_platform_config.set_activated()

        self.logger.info("activation completed")

    def _check_internet_connection(self):
        check_url = 'http://apps.syncloud.org/releases/stable/index'
        internet_ok = True
        try:
            response = requests.get(check_url)
            self.logger.info('Internet check, response status_code: {0}'.format(response.status_code))
            if response.status_code != 200:
                internet_ok = False
         
        except Exception as e:
            self.logger.error('Internet check url {0} is not reachable, error: {1}'.format(check_url, str(e)))
            internet_ok = False
        
        if not internet_ok:
            raise Exception('Internet is not available, check your device connection')

    def set_access(self, upnp_enabled, external_access, manual_public_ip, manual_certificate_port, manual_access_port):
        self.logger.info('set_access: external_access={0}'.format(external_access))

        drill = self.port_drill_factory.get_drill(upnp_enabled, external_access, manual_public_ip,
                                                  manual_certificate_port, manual_access_port)

        if drill is None:
            self.logger.error('Will not change access mode. Was not able to get working port mapper.')
            return

        drill.sync_new_port(WEB_CERTIFICATE_PORT, http_network_protocol)
        mapping = drill.sync_new_port(WEB_ACCESS_PORT, http_network_protocol)
        router_port = None
        if mapping:
            router_port = mapping.external_port

        external_ip = drill.external_ip()

        if self.user_platform_config.is_redirect_enabled():
            self.redirect_service.sync(external_ip, router_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                       self.user_platform_config.get_domain_update_token(), external_access)

        self.user_platform_config.update_device_access(upnp_enabled, external_access,
                                                       manual_public_ip, manual_certificate_port, manual_access_port)
        self.event_trigger.trigger_app_event_domain()

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        port_drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                       manual_certificate_port, manual_access_port)
        try:
            port_drill.sync_existing_ports()
        except Exception as e:
            self.logger.error('Unable to sync port mappings: {0}'.format(str(e)))

        router_port = None
        mapping = port_drill.get(WEB_ACCESS_PORT, http_network_protocol)
        if mapping:
            router_port = mapping.external_port
        
        external_ip = port_drill.external_ip()
        
        self.redirect_service.sync(external_ip, router_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                   update_token, external_access)

    def add_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                  manual_certificate_port, manual_access_port)
        drill.sync_new_port(local_port, protocol)

    def remove_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                  manual_certificate_port, manual_access_port)
        drill.remove(local_port, protocol)


def parse_username(username, domain):
    if '@' in username:
        result = username.split('@')
        name = result[0]
        return name, username
    else:
        email = '{0}@{1}'.format(username, domain)
        return username, email
