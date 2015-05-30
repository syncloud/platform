import requests


def test_main_site_secured():
    response = requests.get('http://localhost/server/rest/user', allow_redirects=False)
    print(response.text)
    assert response.status_code == 302
