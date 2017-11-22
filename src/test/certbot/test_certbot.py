from datetime import datetime

from syncloud_platform.certbot.certbot_generator import apps_to_certbot_domain_args, expiry_date_string_to_days, \
    get_new_domains
from syncloud_platform.sam.models import AppVersions, App


def test_apps_to_certbot_domain_args():

    app_versions = AppVersions()
    app_versions.app = App()
    app_versions.app.id = 'app1'

    domain_args = apps_to_certbot_domain_args([app_versions], 'domain')

    assert domain_args.startswith('-d domain '), 'master domain should be first'


def test_expiry_date_string_to_days_valid():
    assert expiry_date_string_to_days('20171027120200Z', datetime(2017, 10, 20)) == 7


def test_expiry_date_string_to_days_expired():
    assert expiry_date_string_to_days('20171027120200Z', datetime(2017, 10, 31)) == -4


def test_new_domains_more():
    assert get_new_domains(["a", "b"], ["a"]) == ["b"], 'regenerate certificate for new apps'


def test_new_domains_less():
    assert get_new_domains(["a"], ["a", "b"]) == [], 'leave ald apps in the certificate'


def test_new_domains_same():
    assert get_new_domains(["a", "b"], ["a", "b"]) == [], 'do not regenerate with the same list of apps'
