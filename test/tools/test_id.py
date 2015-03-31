from syncloud.tools.footprint import Footprint
from syncloud.tools.id import getname


def test_getname_cubietruck():
    f = Footprint('sun7i', cpu_count=2, mem_size=1911201792)
    name = getname(f)
    assert name == 'cubietruck'


def test_getname_unknown():
    f = Footprint('i386', cpu_count=16, mem_size=2000000000)
    name = getname(f)
    assert name is None