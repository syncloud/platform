import subprocess
from os.path import join
from subprocess import check_output

from datetime import datetime

from os import path
from syncloudlib.logger import get_logger

from syncloud_platform.certificate.certbot.certbot_result import CertbotResult
from OpenSSL import crypto
from pyasn1.codec.der import decoder as der_decoder
from ndg.httpsclient.ssl_peer_verification import SUBJ_ALT_NAME_SUPPORT
from ndg.httpsclient.subj_alt_name import SubjectAltName


def apps_to_certbot_domain_args(app_versions, domain):
    domains = domain_list_sorted(app_versions, domain)
    domain_args = '-d ' + ' -d '.join(domains)
    return domain_args


def domain_list_sorted(app_versions, domain):
    # we need to list all the individual domains for now as wildcard domain is not supported by certbot yet
    all_apps = [app_versions.app.id for app_versions in app_versions]
    domains = ['{0}.{1}'.format(app, domain) for app in all_apps]
    domains.sort()
    domains.insert(0, domain)
    return domains


class CertbotGenerator:
    def __init__(self, platform_config, user_platform_config, info, snap):
        self.info = info
        self.log = get_logger('certbot')
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.certbot_bin = '{0}/bin/certbot'.format(self.platform_config.app_dir())
        self.log_dir = self.platform_config.get_log_root()
        self.certbot_config_dir = join(self.platform_config.data_dir(), 'certbot')
        self.snap = snap

    def certbot_certificate_file(self):
        return '{0}/certbot/live/{1}/fullchain.pem'.format(
            self.platform_config.data_dir(), self.info.domain())

    def certbot_key_file(self):
        return '{0}/certbot/live/{1}/privkey.pem'.format(
                self.platform_config.data_dir(), self.info.domain())

    def generate_certificate(self, is_test_cert=False):

        self.log.info('running certbot')
        domain_args = apps_to_certbot_domain_args(self.snap.list(), self.info.domain())

        test_cert = ''
        if is_test_cert:
            test_cert = '--test-cert --break-my-certs'

        plugin = '--webroot --webroot-path {0}/certbot/www'.format(self.platform_config.data_dir())

        try:
            cmd =  '{0} --logs-dir={1} --max-log-backups 5 --config-dir={2} --agree-tos --email {3} certonly --force-renewal --cert-name {4} {5} {6} {7} '.format(
                    self.certbot_bin, self.log_dir, self.certbot_config_dir,
                    self.user_platform_config.get_user_email(), self.info.domain(),
                    test_cert, plugin, domain_args
                    )
            self.log.info(cmd)
            output = check_output(cmd, stderr=subprocess.STDOUT, shell=True).decode()
            self.log.info(output)
            archive_dir = join(self.certbot_config_dir, 'archive')
            if path.exists(archive_dir):
                check_output('chmod 755 {0}'.format(archive_dir), shell=True)
            live_dir = join(self.certbot_config_dir, 'live')
            if path.exists(live_dir):
                check_output('chmod 755 {0}'.format(live_dir), shell=True)
            if not path.exists(self.certbot_certificate_file()):
                raise Exception("certificate does not exist: {0}".format(self.certbot_certificate_file()))
            return CertbotResult(self.certbot_certificate_file(), self.certbot_key_file())

        except subprocess.CalledProcessError as e:
            self.log.warn(e.output.decode())
            raise e

    def days_until_expiry(self):

        self.log.info('getting expiry date of {}'.format(self.certbot_certificate_file()))
        if not path.exists(self.certbot_certificate_file()):
            self.log.info('certificate does not exist yet, {0}'.format(self.certbot_certificate_file()))
            return 0

        cert = crypto.load_certificate(crypto.FILETYPE_PEM, file(self.certbot_certificate_file()).read())
        days = expiry_date_string_to_days(cert.get_notAfter())
        return days

    def new_domains(self):

        current_domains = domain_list_sorted(self.snap.list(), self.info.domain())

        cert_domains = []
        if path.isfile(self.certbot_certificate_file()):
            cert = crypto.load_certificate(crypto.FILETYPE_PEM, file(self.certbot_certificate_file()).read())
            cert_domains = get_subj_alt_name(cert)

        return get_new_domains(current_domains, cert_domains)


# Note: This is a slightly bug-fixed version of same from ndg-httpsclient.
def get_subj_alt_name(peer_cert):
    # Search through extensions
    dns_name = []
    if not SUBJ_ALT_NAME_SUPPORT:
        return dns_name

    general_names = SubjectAltName()
    for i in range(peer_cert.get_extension_count()):
        ext = peer_cert.get_extension(i)
        ext_name = ext.get_short_name()
        if ext_name != 'subjectAltName':
            continue

        # PyOpenSSL returns extension data in ASN.1 encoded form
        ext_dat = ext.get_data()
        decoded_dat = der_decoder.decode(ext_dat,
                                         asn1Spec=general_names)

        for name in decoded_dat:
            if not isinstance(name, SubjectAltName):
                continue
            for entry in range(len(name)):
                component = name.getComponentByPosition(entry)
                if component.getName() != 'dNSName':
                    continue
                dns_name.append(str(component.getComponent()))

    return dns_name


def expiry_date_string_to_days(expiry, today=datetime.today()):
    expiry_date = datetime.strptime(expiry, "%Y%m%d%H%M%SZ")
    return (expiry_date - today).days


def get_new_domains(current_domains, cert_domains):
    return list(set(current_domains) - set(cert_domains))

