from syncloud_platform.control.tls import apps_to_certbot_domain_args
from syncloud_platform.sam.models import AppVersions, App


def test_apps_to_certbot_domain_args():

    app_versions = AppVersions()
    app_versions.app = App()
    app_versions.app.id = 'app1'

    domain_args = apps_to_certbot_domain_args([app_versions], 'domain')

    assert domain_args.startswith('-d domain '), 'master domain should be first'
