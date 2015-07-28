import syncloud_platform.importlib

from syncloud_platform.tools.config import footprints, titles

def test_all_known_footprints_are_different():
    for name, f in footprints:
        print name
        print f
        same = {name2: f2 for (name2, f2) in footprints if f == f2}
        assert len(same) == 1, 'These footprints are intersecting: {}'.format(', '.join(set(same.keys())))

def test_all_names_have_titles():
    no_titles = [name for name, f in footprints if name not in titles]
    assert len(no_titles) == 0, 'These names do not have titles: {}'.format(', '.join(no_titles))