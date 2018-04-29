import logging

from syncloud_platform.snap.snap import Snap
from syncloud_app import logger

snaps_json = """{"type":"sync","status-code":200,"status":"OK","result":[{"apps":[],"broken":"","channel":"stable","confinement":"strict","description":"The core runtime environment for snapd","developer":"canonical","devmode":false,"icon":"","id":"99T7MUlRhtI3U0QFgl5mXXESAiSwt776","install-date":"2016-12-16T04:28:38Z","installed-size":79720448,"jailmode":false,"name":"core","private":false,"resource":"/v2/snaps/core","revision":"714","status":"active","summary":"snapd runtime environment","trymode":false,"type":"os","version":"16.04.1"},{"apps":[{"name":"openldap","daemon":"forking","aliases":null},{"name":"uwsgi-internal","daemon":"notify","aliases":null},{"name":"uwsgi-public","daemon":"notify","aliases":null},{"name":"nginx","daemon":"forking","aliases":null}],"broken":"","channel":"","confinement":"strict","description":"Syncloud service store.","developer":"","devmode":true,"icon":"","id":"","install-date":"2017-01-11T18:57:24Z","installed-size":35192832,"jailmode":false,"name":"syncloud-platform","private":false,"resource":"/v2/snaps/syncloud-platform","revision":"x1","status":"active","summary":"Syncloud Platform","trymode":false,"type":"app","version":"1188"}],"sources":["local"]}"""


class Info:
    def url(self, id):
        return ""


logger.init(logging.DEBUG, True)


def test_parse_snaps_response():
    apps = Snap(None, Info()).parse_response(snaps_json, lambda app: True)
    assert len(apps) == 2
