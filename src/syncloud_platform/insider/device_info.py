class DeviceInfo:
    def __init__(self, user_platform_config):
        self.user_platform_config = user_platform_config

    def domain(self):
        domain = 'localhost'

        if self.user_platform_config.is_redirect_enabled():
            full_domain = self.user_platform_config.get_domain()
            if full_domain is not None:
                domain = full_domain
            else:
                user_domain = self.user_platform_config.get_user_domain()
                if user_domain is not None:
                    domain = '{0}.{1}'.format(user_domain, self.user_platform_config.get_redirect_domain())
        else:
            custom = self.user_platform_config.get_custom_domain()
            if custom:
                domain = custom
        return domain
