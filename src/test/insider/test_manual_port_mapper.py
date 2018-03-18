from syncloud_platform.insider.manual import ManualPortMapper


def test_add_mapping_manual_web():
    mapper = ManualPortMapper('192.168.0.1', 8080, 8081)
    assert mapper.add_mapping(80, 80, 'TCP') == 8080
    assert mapper.add_mapping(443, 81, 'TCP') == 8081

def test_add_mapping_manual_non_web():
    mapper = ManualPortMapper('192.168.0.1', 8080, 8081)
    assert mapper.add_mapping(25, 25, 'TCP') == 25

