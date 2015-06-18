from syncloud.tools.nginx import Nginx


def test_proxy():
    assert Nginx().proxy_definition('test', 80) == 'location /test {\n' \
                                                   '    proxy_pass      http://localhost:80/test/ ;\n' \
                                                   '    proxy_redirect  http://localhost:80/test/ /test/ ;\n' \
                                                   '}'
