from syncloud.tools.nginx import Nginx


def test_proxy():
    assert Nginx().proxy_definition('test', 80) == 'location /test { ' \
                                                   '   proxy_pass http://127.0.0.1:80; ' \
                                                   '   proxy_redirect  http://localhost:80/ /;' \
                                                   '}'
