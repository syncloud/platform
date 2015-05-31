import requests
from subprocess import check_output
from socket import gethostname;

def test_avahi():
    discovery = check_output('/opt/app/platform/avahi/bin/avahi-browse -atk', shell=True)
    assert 'syncloud on {0}'.format(gethostname()) in discovery


def test_main_site_secured():

    response = requests.get('http://localhost/server/rest/user', allow_redirects=False)
    print(response.text)
    assert response.status_code == 302
