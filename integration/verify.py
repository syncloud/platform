import requests
from subprocess import check_output


def test_avahi():
    assert '' in check_output('/opt/app/platform/avahi/bin/avahi-browse -atk', shell=True)


def test_main_site_secured():

    response = requests.get('http://localhost/server/rest/user', allow_redirects=False)
    print(response.text)
    assert response.status_code == 302
