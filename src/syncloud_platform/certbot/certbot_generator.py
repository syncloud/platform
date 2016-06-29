import subprocess
from os.path import join
from subprocess import check_output

from syncloud_app.logger import get_logger

from syncloud_platform.certbot.certbot_result import CertbotResult


def apps_to_certbot_domain_args(app_versions, domain):
    # we need to list all the individual domains for now as wildcard domain is not supported by certbot yet
    all_apps = [app_versions.app.id for app_versions in app_versions]
    domains = ['{0}.{1}'.format(app, domain) for app in all_apps]
    domains.append(domain)
    domains.reverse()
    domain_args = '-d ' + ' -d '.join(domains)
    return domain_args


class CertbotGenerator:
    def __init__(self, platform_config, user_platform_config, info, sam):
        self.info = info
        self.log = get_logger('certbot')
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.certbot_bin = '{0}/bin/certbot'.format(self.platform_config.app_dir())
        self.log_dir = self.platform_config.get_log_root()
        self.certbot_config_dir = join(self.platform_config.data_dir(), 'certbot')
        self.sam = sam
        self.certbot_certificate_file = '{0}/certbot/live/{1}/fullchain.pem'.format(
            self.platform_config.data_dir(), self.info.domain())
        self.certbot_key_file = '{0}/certbot/keys/0000_key-certbot.pem'.format(
                self.platform_config.data_dir(), self.info.domain())

    def generate_certificate(self, is_test_cert):

        self.log.info('running certbot')
        domain_args = apps_to_certbot_domain_args(self.sam.list(), self.info.domain())

        test_cert = ''
        if is_test_cert:
            test_cert = '--test-cert'

        output = check_output(
            '{0} --logs-dir={1} --config-dir={2} --agree-tos --email {3} '
            'certonly {4} --webroot --webroot-path {5} '
            '{6} '.format(self.certbot_bin,
                          self.log_dir,
                          self.certbot_config_dir,
                          self.user_platform_config.get_user_email(),
                          test_cert,
                          self.platform_config.www_root(),
                          domain_args), stderr=subprocess.STDOUT, shell=True)

        self.log.info(output)

        regenerated = 'no action taken' not in output
               
        return CertbotResult(self.certbot_certificate_file, self.certbot_key_file, regenerated)
