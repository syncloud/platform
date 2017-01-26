from syncloud_platform.board.id import id

def test_unknown():
    the_id = id('/etc/non-existing-file.cfg')
    assert the_id.name == 'unknown'
    assert the_id.title == 'Unknown'