import logging
from os.path import dirname
from syncloud.app import logger
from syncloud.tools.cpu.cpuinfo import CpuInfo
from syncloud.tools.cpu.reader import Reader

test_dir = dirname(__file__)
logger.init(logging.DEBUG, console=True)


def test_hardware():

    assert CpuInfo(Reader(device('beagle-bone-black'))).hardware() == 'Generic AM33XX (Flattened Device Tree)'
    assert CpuInfo(Reader(device('cubieboard'))).hardware() == 'sun4i'
    assert CpuInfo(Reader(device('cubieboard2'))).hardware() == 'sun7i'
    assert CpuInfo(Reader(device('cubietruck'))).hardware() == 'sun7i'
    assert CpuInfo(Reader(device('raspberry-pi-b'))).hardware() == 'BCM2708'


def device(device_file):
    return test_dir + "/cpu/" + device_file
