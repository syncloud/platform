from os.path import dirname, join

from syncloud_platform.tools.hardware import Hardware

DIR = dirname(__file__)


def test_list():
    disks = Hardware(open(join(DIR, 'hardware', 'lshw.json')).read()).disks()
    assert len(disks) == 3
