from os.path import dirname, join

import syncloud_platform.importlib

from syncloud_platform.tools.nginx import Nginx


dir = dirname(__file__)


def test_proxy_location():
    assert Nginx().proxy_definition('test', 80, join(dir, '..', '..', '..', 'config', 'nginx'), 'app.location') == """
location /test {
    proxy_set_header X-Forwarded-Proto $scheme ;
    proxy_set_header X-Forwarded-Host $http_host ;
    proxy_pass      http://localhost:80/test ;
    proxy_redirect  http://localhost:80/test $scheme://$http_host/test ;
}
""".strip()


def test_proxy_server():
    assert Nginx().proxy_definition('test', 80, join(dir, '..', '..', '..', 'config', 'nginx'), 'app.server') == """
server {
    listen 80;
    server_name test.*;
    location / {
        proxy_set_header X-Forwarded-Proto $scheme ;
        proxy_set_header X-Forwarded-Host $http_host ;
        proxy_pass      http://localhost:80 ;
        proxy_redirect  http://localhost:80 $scheme://$http_host ;
    }
}

server {

    listen 443 ssl;
    server_name test.*;

    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains";

    location / {
        proxy_set_header X-Forwarded-Proto $scheme ;
        proxy_set_header X-Forwarded-Host $http_host ;
        proxy_pass      http://localhost:80 ;
        proxy_redirect  http://localhost:80 $scheme://$http_host ;
    }
}

""".strip()
