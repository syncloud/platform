from syncloud.insider.config import Service
from syncloud.insider.dns import Endpoint
from syncloud.server.model import Site


def test_site_from_endpoint():
    endpoint = Endpoint(Service("image-ci", "http", "type", "80", "image-ci-url"), 'localhost', 8181)
    site = Site(endpoint)
    assert site.name == 'image-ci'
    assert site.url == 'http://localhost:8181/image-ci-url'