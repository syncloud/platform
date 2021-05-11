package activation

import "fmt"

type Free struct {
}

func New() *Free {
	return &Free{}
}

func (f *Free) Activate() error {
	/*
		def activate(self, redirect_email, redirect_password, user_domain, device_username, device_password, main_domain):
		        user_domain_lower = user_domain.lower()
		        self.logger.info("activate {0}, {1}".format(user_domain_lower, device_username))

		        self._check_internet_connection()

		        user = self.prepare_redirect(redirect_email, redirect_password, main_domain)
		        self.user_platform_config.set_user_update_token(user.update_token)

		        name, email = parse_username(device_username, '{0}.{1}'.format(user_domain_lower, main_domain))

		        response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain_lower)
		        self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)

		        self._activate_common(name, device_username, device_password, email)
	*/
	return fmt.Errorf("not implemented yet")
}
