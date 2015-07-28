import syncloud_platform.importlib

from syncloud_platform.tools.footprint import Footprint


def test_equals_all_members():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('arm', cpu_count=1, mem_size=1234)
    assert footprint.match(pattern)


def test_equals_no_cpu_count():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('arm', mem_size=1234)
    assert footprint.match(pattern)


def test_equals_different_cpu_count():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('arm', cpu_count=2, mem_size=1234)
    assert not footprint.match(pattern)


def test_equals_no_mem_size():
    footprint = Footprint('arm', cpu_count=1)
    pattern = Footprint('arm', cpu_count=1, mem_size=1234)
    assert not footprint.match(pattern)


def test_equals_different_mem_size():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('arm', cpu_count=1, mem_size=5678)
    assert not footprint.match(pattern)


def test_equals_different_cpu_hardware():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('i386', cpu_count=1, mem_size=1234)
    assert not footprint.match(pattern)


def test_equals_same_cpu_hardware():
    footprint = Footprint('arm', cpu_count=1, mem_size=1234)
    pattern = Footprint('arm', cpu_count=1, mem_size=1234)
    assert footprint.match(pattern)


def test_equals_vendor():
    footprint = Footprint(vendor_id='AuthenticAMD', cpu_count=32, mem_size=1234)
    pattern = Footprint(vendor_id='AuthenticAMD')
    assert footprint.match(pattern)