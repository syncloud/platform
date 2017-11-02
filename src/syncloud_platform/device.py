import getpass
import uuid

from syncloud_app import logger

from syncloud_platform.insider.util import protocol_to_port, secure_to_protocol
from syncloud_platform.gaplib import fs

http_network_protocol = 'TCP'


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service,
                 port_drill_factory, sam, platform_cron, ldap_auth, event_trigger, tls, nginx):
        self.tls = tls
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory
        self.sam = sam
        self.auth = ldap_auth
        self.platform_cron = platform_cron
        self.event_trigger = event_trigger
        self.logger = logger.get_logger('Device')
        self.nginx = nginx

    def prepare_redirect(self, redirect_email, redirect_password, main_domain):

        redirect_api_url = 'http://api.' + main_domain

        self.logger.info("prepare redirect {0}, {1}".format(redirect_email, redirect_api_url))
        self.user_platform_config.set_redirect_enabled(True)
        self.sam.update()
        self.user_platform_config.update_redirect(main_domain, redirect_api_url)
        self.user_platform_config.set_user_email(redirect_email)

        user = self.redirect_service.get_user(redirect_email, redirect_password)
        return user

    def activate(self, redirect_email, redirect_password, user_domain, device_username, device_password, main_domain):

        self.logger.info("activate {0}, {1}".format(user_domain, device_username))

        user = self.prepare_redirect(redirect_email, redirect_password, main_domain)
        self.user_platform_config.set_user_update_token(user.update_token)

        response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain)
        self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)

        self.platform_cron.remove()
        self.platform_cron.create()

        self.set_access(False, False, False, 0, 0)

        self.logger.info("activating ldap")
        fix_permissions = self.platform_config.get_installer() == 'sam'
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        self.tls.generate_self_signed_certificate()
        user, email = parse_username(device_username, '{0}.{1}'.format(user_domain, main_domain))
        self.auth.reset(user, device_username, device_password, fix_permissions, email)
        
        self.nginx.init_config()
        self.nginx.reload_public()
        
        self.logger.info("activation completed")

    def activate_custom_domain(self, full_domain, device_username, device_password):

        self.logger.info("activate custom {0}, {1}".format(full_domain, device_username))
        self.sam.update()
        
        self.user_platform_config.set_redirect_enabled(False)
        self.user_platform_config.set_custom_domain(full_domain)
        
        user, email = parse_username(device_username, full_domain)
        self.user_platform_config.set_user_email(email)

        self.platform_cron.remove()
        self.platform_cron.create()

        self.set_access(False, False, False, 0, 0)

        self.logger.info("activating ldap")
        fix_permissions = self.platform_config.get_installer() == 'sam'
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        self.tls.generate_self_signed_certificate()

        self.auth.reset(user, device_username, device_password, fix_permissions, email)
        
        self.nginx.init_config()
        self.nginx.reload_public()
        
        self.logger.info("activation completed")

    def set_access(self, upnp_enabled, is_https, external_access, manual_public_ip, manual_public_port):
        self.logger.info('set_access: https={0}, external_access={1}'.format(is_https, external_access))

        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        web_protocol = secure_to_protocol(is_https)
        new_web_local_port = protocol_to_port(web_protocol)
        old_web_protocol = secure_to_protocol(self.user_platform_config.is_https())
        old_web_local_port = protocol_to_port(old_web_protocol)

        drill = self.port_drill_factory.get_drill(upnp_enabled, external_access, manual_public_ip, manual_public_port)

        if drill is None:
            self.logger.error('Will not change access mode. Was not able to get working port mapper.')
            return
        try:
            drill.remove(old_web_local_port, http_network_protocol)
        except Exception, e:
            self.logger.error('Unable to remove port {0}: {1}'.format(old_web_local_port, e.message))

        drill.sync_new_port(new_web_local_port, http_network_protocol)

        self.redirect_service.sync(drill, update_token, web_protocol, external_access, http_network_protocol)
        self.user_platform_config.update_device_access(upnp_enabled, is_https, external_access,
                                                       manual_public_ip, manual_public_port)
        self.event_trigger.trigger_app_event_domain()

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_public_port = self.user_platform_config.get_manual_public_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip, manual_public_port)
        web_protocol = secure_to_protocol(self.user_platform_config.is_https())
        self.redirect_service.sync(drill, update_token, web_protocol, external_access, http_network_protocol)

        if not getpass.getuser() == self.platform_config.cron_user():
            fs.chownpath(self.platform_config.data_dir(), self.platform_config.cron_user())

    def add_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_public_port = self.user_platform_config.get_manual_public_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip, manual_public_port)
        drill.sync_new_port(local_port, protocol)

    def remove_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_public_port = self.user_platform_config.get_manual_public_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip, manual_public_port)
        drill.remove(local_port, protocol)


def parse_username(username, domain):
    if '@' in username:
        result = username.split('@')
        name = result[0]
        return name, username
    else:
        email = '{0}@{1}'.format(username, domain)
        return username, email

        